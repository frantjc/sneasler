package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/frantjc/go-ingress"
	"github.com/frantjc/sneasler"
	wellknown "github.com/frantjc/sneasler/.well-known"
	wellknownblob "github.com/frantjc/sneasler/.well-known/blob"
	"github.com/frantjc/sneasler/api"
	"github.com/frantjc/sneasler/internal/conf"
	"github.com/frantjc/sneasler/internal/openapi"
	"github.com/frantjc/sneasler/iofunc"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/logsquaredn/blobproxy/bucketfs"
	"github.com/spf13/cobra"
	"gocloud.dev/blob"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func NewSneasler() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "sneasler",
			Version:       sneasler.GetSemver(),
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, _ []string) {
				cmd.SetContext(
					sneasler.WithLogger(cmd.Context(), sneasler.NewLogger().V(2-verbosity)),
				)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx             = cmd.Context()
					log             = sneasler.LoggerFrom(ctx)
					errC            = make(chan error, 1)
					l               net.Listener
					title           = "Sneasler"
					frontendHandler = http.NotFoundHandler()
				)

				if conf.Viper.IsSet("js-entrypoint") {
					// try to cleverly pick a random unused port to bind frontend to
					rl, err := net.Listen("tcp", "127.0.0.1:0")
					if err != nil {
						return err
					}

					var (
						_, frontendPort, _ = net.SplitHostPort(rl.Addr().String())
						frontendEntrypoint = conf.Viper.GetString("js-entrypoint")
						frontendWd         = filepath.Dir(frontendEntrypoint)
						//nolint:gosec
						frontend       = exec.CommandContext(ctx, conf.Viper.GetString("node"), frontendEntrypoint)
						frontendURL, _ = url.Parse("http://localhost:" + frontendPort + "/")
						frontendLog    = log.WithName("frontend")
						stdoutErr      = errors.New("stderr log")
					)
					
					if err := rl.Close(); err != nil {
						return err
					}

					go func() {
						frontend.Dir = frontendWd
						frontend.Env = append(frontend.Environ(), "PORT="+frontendPort, "NODE_ENV=production")
						frontend.Stdout = iofunc.WriterFunc(func(b []byte) (int, error) {
							frontendLog.Info(strings.TrimSpace(string(b)))
							return len(b), nil
						})
						frontend.Stderr = iofunc.WriterFunc(func(b []byte) (int, error) {
							frontendLog.Error(stdoutErr, strings.TrimSpace(string(b)))
							return len(b), nil
						})
						errC <- frontend.Run()
					}()

					frontendHandler = httputil.NewSingleHostReverseProxy(frontendURL)
				}

				openapi.Spec.Info = &spec.Info{
					InfoProps: spec.InfoProps{
						Title:   title,
						Version: "v" + cmd.Version,
					},
				}

				addr, err := url.Parse(args[0])
				if err != nil {
					return err
				}

				bucket, err := blob.OpenBucket(ctx, addr.String())
				if err != nil {
					return err
				}
				defer bucket.Close()

				if accessible, err := bucket.IsAccessible(ctx); !accessible || err != nil {
					return fmt.Errorf("inaccessible bucket " + addr.String())
				}

				if conf.Viper.GetBool("ngrok") {
					tun, err := ngrok.Listen(
						ctx,
						config.HTTPEndpoint(),
						ngrok.WithAuthtokenFromEnv(),
					)
					if err != nil {
						return err
					}
					defer tun.CloseWithContext(ctx) //nolint:errcheck

					log.Info("listening", "url", tun.URL())
					u, _ := url.Parse(tun.URL())
					openapi.Spec.Host = u.Host
					openapi.Spec.Info.Contact = &spec.ContactInfo{
						ContactInfoProps: spec.ContactInfoProps{
							URL: u.String(),
						},
					}
					l = tun
				} else {
					p, err := net.Listen("tcp", fmt.Sprint(":", conf.Viper.GetInt64("port")))
					if err != nil {
						return err
					}
					defer p.Close()

					log.Info("listening", "addr", p.Addr())
					l = p
				}

				var (
					b, _     = json.MarshalIndent(openapi.Spec, "", "  ")
					lenB     = len(b)
					specURL  = "/swagger.json"
					docsPath = "/docs"
					next     = http.NotFoundHandler()
					rapiDoc  = middleware.RapiDoc(
						middleware.RapiDocOpts{
							Title:   title,
							Path:    docsPath,
							SpecURL: specURL,
						},
						next,
					)
					redoc = middleware.Redoc(
						middleware.RedocOpts{
							Title:   title,
							Path:    docsPath,
							SpecURL: specURL,
						},
						next,
					)
					swaggerUI = middleware.SwaggerUI(
						middleware.SwaggerUIOpts{
							Title:   title,
							Path:    docsPath,
							SpecURL: specURL,
						},
						next,
					)
					srv = http.Server{
						ReadHeaderTimeout: 30 * time.Second,
						ReadTimeout:       30 * time.Second,
						WriteTimeout:      30 * time.Second,
						IdleTimeout:       2 * time.Minute,
						BaseContext: func(_ net.Listener) context.Context {
							return ctx
						},
						Handler: ingress.New(
							ingress.PrefixPath(
								wellknown.Path,
								bucketfs.NewFileServer(bucketfs.NewFS(bucket).WithContext(ctx)),
							),
							ingress.PrefixPath(
								"/api",
								api.NewHandler(
									&wellknownblob.AssetlinksBucket{
										Bucket: bucket,
										Key:    wellknown.AssetlinksPath,
									},
									&wellknownblob.AppleAppSiteAssociationBucket{
										Bucket: bucket,
										Key:    wellknown.AppleAppSiteAssociationPath,
									},
								),
							),
							ingress.PrefixPath(
								specURL,
								http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									if r.URL.Path == specURL {
										w.Header().Set("Content-Type", "application/json")
										w.Header().Set("Content-Length", fmt.Sprint(lenB))
										_, _ = w.Write(b)
										return
									}

									next.ServeHTTP(w, r)
								}),
							),
							ingress.PrefixPath(
								docsPath,
								// RapiDoc, SwaggerUI, Redoc are equivalent options here
								http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									switch strings.ToLower(r.URL.Query().Get("flavor")) {
									case "", "rapidoc":
										rapiDoc.ServeHTTP(w, r)
									case "redoc":
										redoc.ServeHTTP(w, r)
									case "swagger":
										swaggerUI.ServeHTTP(w, r)
									}
								}),
							),
							ingress.PrefixPath("/", frontendHandler),
						),
					}
				)

				go func() {
					srv.SetKeepAlivesEnabled(true)
					errC <- srv.Serve(l)
				}()

				select {
				case <-ctx.Done():
					return errors.Join(ctx.Err(), srv.Shutdown(ctx))
				case err := <-errC:
					return err
				}
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")

	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity")

	cmd.Flags().Int64P("port", "p", mustParsePort(), "port")
	_ = conf.Viper.BindPFlag("port", cmd.Flag("port"))
	conf.Viper.MustBindEnv("port", "PORT")

	cmd.Flags().Bool("ngrok", false, "use ngrok")
	_ = conf.Viper.BindPFlag("ngrok", cmd.Flag("ngrok"))
	conf.Viper.MustBindEnv("ngrok", "SNEASLER_NGROK")

	cmd.Flags().String("node", "node", "NodeJS")
	_ = conf.Viper.BindPFlag("node", cmd.Flag("node"))
	conf.Viper.MustBindEnv("node", "SNEASLER_NODE")

	cmd.Flags().String("js-entrypoint", "", "JavaScript entrypoint")
	_ = conf.Viper.BindPFlag("js-entrypoint", cmd.Flag("js-entrypoint"))
	conf.Viper.MustBindEnv("js-entrypoint", "SNEASLER_JS_ENTRYPOINT")

	return cmd
}

func mustParsePort() int64 {
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if p != 0 && err != nil {
		return int64(p)
	}

	return 8080
}

package command

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/frantjc/go-ingress"
	"github.com/frantjc/sneasler"
	wellknown "github.com/frantjc/sneasler/.well-known"
	wellknownblob "github.com/frantjc/sneasler/.well-known/blob"
	"github.com/frantjc/sneasler/api"
	"github.com/logsquaredn/blobproxy/bucketfs"
	"github.com/spf13/cobra"
	"gocloud.dev/blob"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func NewSneasler() *cobra.Command {
	var (
		verbosity int
		port      int64
		useNgrok  bool
		cmd       = &cobra.Command{
			Use:           "sneasler",
			Version:       sneasler.GetSemver(),
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					sneasler.WithLogger(cmd.Context(), sneasler.NewLogger().V(2-verbosity)),
				)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx  = cmd.Context()
					log  = sneasler.LoggerFrom(ctx)
					errC = make(chan error, 1)
					l    net.Listener
				)

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
					return fmt.Errorf("inaccessible bucket %s", addr.String())
				}

				if useNgrok {
					tun, err := ngrok.Listen(
						ctx,
						config.HTTPEndpoint(),
						ngrok.WithAuthtokenFromEnv(),
					)
					if err != nil {
						return err
					}
					defer tun.CloseWithContext(ctx)

					log.Info("listening", "url", tun.URL())
					l = tun
				} else {
					p, err := net.Listen("tcp", fmt.Sprint(":", port))
					if err != nil {
						return err
					}
					defer p.Close()

					log.Info("listening", "addr", p.Addr())
					l = p
				}

				var (
					aasab = &wellknownblob.AppleAppSiteAssociationBucket{
						Bucket: bucket,
						Key:    wellknown.AppleAppSiteAssociationPath,
					}
					alb = &wellknownblob.AssetlinksBucket{
						Bucket: bucket,
						Key:    wellknown.AssetlinksPath,
					}
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
								api.NewHandler(alb, aasab),
							),
						),
					}
				)

				go func() {
					errC <- srv.Serve(l)
				}()
				defer srv.Shutdown(ctx) //nolint:errcheck

				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-errC:
					return err
				}
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity")
	cmd.Flags().Int64VarP(&port, "port", "p", mustParsePort(), "port")
	cmd.Flags().BoolVar(&useNgrok, "ngrok", false, "ngrok")

	return cmd
}

func mustParsePort() int64 {
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if p != 0 && err != nil {
		return int64(p)
	}

	return 8080
}

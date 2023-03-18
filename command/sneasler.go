package command

import (
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/frantjc/sneasler"
	"github.com/spf13/cobra"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func NewSneasler() *cobra.Command {
	var (
		verbosity int
		port      int64
		cmd       = &cobra.Command{
			Use:           "sneasler",
			Version:       sneasler.GetSemver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					sneasler.WithLogger(cmd.Context(), sneasler.NewLogger().V(2-verbosity)),
				)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx = cmd.Context()
					log = sneasler.LoggerFrom(ctx)
				)

				tun, err := ngrok.Listen(ctx,
					config.HTTPEndpoint(),
					ngrok.WithAuthtokenFromEnv(),
				)
				if err != nil {
					return err
				}

				http.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("hello"))
				})

				log.Info("listening", "url", tun.URL())
				return http.Serve(tun, nil) //nolint:gosec
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity")
	cmd.Flags().Int64VarP(&port, "port", "p", mustParsePort(), "port")

	return cmd
}

func mustParsePort() int64 {
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if p != 0 && err != nil {
		return int64(p)
	}

	return 8080
}

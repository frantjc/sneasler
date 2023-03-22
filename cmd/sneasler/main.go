package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	errorcode "github.com/frantjc/go-error-code"
	"github.com/frantjc/sneasler/command"
	"github.com/gin-gonic/gin"

	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	var (
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		err       error
	)

	if err = command.NewSneasler().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}

	stop()
	os.Exit(errorcode.ExitCode(err))
}

package api

import (
	"context"
	"net/http"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"github.com/gin-gonic/gin"
)

func assetlinkHandler(alb wellknown.AssetlinksBackend, f func(context.Context, wellknown.AssetlinksBackend, *wellknown.Assetlink) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := &wellknown.Assetlink{}

		if err := ctx.Bind(body); err != nil {
			return
		}

		if err := f(ctx, alb, body); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

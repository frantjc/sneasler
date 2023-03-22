package api

import (
	"context"
	"net/http"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"github.com/gin-gonic/gin"
)

func appleAppSiteAssociationApplinkDetailHandler(aasab wellknown.AppleAppSiteAssociationBackend, f func(context.Context, wellknown.AppleAppSiteAssociationBackend, *wellknown.AppleAppSiteAssociationAppLinksDetail) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := &wellknown.AppleAppSiteAssociationAppLinksDetail{}

		if err := ctx.Bind(body); err != nil {
			return
		}

		if err := f(ctx, aasab, body); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

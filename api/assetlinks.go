package api

import (
	"context"
	"net/http"
	"reflect"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"github.com/frantjc/sneasler/internal/openapi"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
)

func assetlinksHandler(alb wellknown.AssetlinksBackend, f func(context.Context, wellknown.AssetlinksBackend, *wellknown.Assetlink) error) gin.HandlerFunc {
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

func init() {
	openapi.Spec.Paths.Paths["/api/v1/assetlinks"] = spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Post:  assetlinksOperation(),
			Put:   assetlinksOperation(),
			Patch: assetlinksOperation(),
		},
	}

	openapi.Spec.Definitions["assetlink"] = *openapi.Definition(reflect.TypeOf(&wellknown.Assetlink{}))
}

func assetlinksOperation() *spec.Operation {
	return &spec.Operation{
		OperationProps: spec.OperationProps{
			Tags:      []string{"assetlinks"},
			Responses: responses(http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError),
		},
	}
}

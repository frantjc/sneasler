package api

import (
	"context"
	"net/http"

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

	openapi.Spec.Definitions["assetlink"] = *spec.MapProperty(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Properties: spec.SchemaProperties{
				"relation": *spec.ArrayProperty(spec.StringProperty()),
				"target": *spec.MapProperty(&spec.Schema{
					SchemaProps: spec.SchemaProps{
						Properties: spec.SchemaProperties{
							"namespace":                *spec.StringProperty(),
							"package_name":             *spec.StringProperty(),
							"sha256_cert_fingerpritns": *spec.ArrayProperty(spec.StringProperty()),
						},
					},
				}),
			},
		},
	})
}

func assetlinksOperation() *spec.Operation {
	return &spec.Operation{
		OperationProps: spec.OperationProps{
			Tags: []string{
				"assetlinks",
			},
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusOK: {
							ResponseProps: spec.ResponseProps{
								Description: http.StatusText(http.StatusOK),
							},
						},
						http.StatusBadRequest: {
							ResponseProps: spec.ResponseProps{
								Description: http.StatusText(http.StatusBadRequest),
							},
						},
						http.StatusInternalServerError: {
							ResponseProps: spec.ResponseProps{
								Description: http.StatusText(http.StatusInternalServerError),
							},
						},
					},
				},
			},
		},
	}
}

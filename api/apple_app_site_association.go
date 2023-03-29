package api

import (
	"context"
	"net/http"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"github.com/frantjc/sneasler/internal/openapi"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
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

func init() {
	openapi.Spec.Paths.Paths["/api/v1/apple-app-site-association"] = spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Post:  appleAppSiteAssociationOperation(),
			Put:   appleAppSiteAssociationOperation(),
			Patch: appleAppSiteAssociationOperation(),
		},
	}

	openapi.Spec.Definitions["apple-app-site-association.applinks.detail"] = *spec.MapProperty(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Properties: spec.SchemaProperties{
				"appID": *spec.StringProperty(),
				"paths": *spec.ArrayProperty(spec.StringProperty()),
			},
		},
	})
}

func appleAppSiteAssociationOperation() *spec.Operation {
	return &spec.Operation{
		OperationProps: spec.OperationProps{
			Tags: []string{
				"apple-app-site-association",
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

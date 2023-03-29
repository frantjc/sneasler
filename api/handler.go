package api

import (
	"net/http"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"github.com/gin-gonic/gin"
)

func NewHandler(alb wellknown.AssetlinksBackend, aasab wellknown.AppleAppSiteAssociationBackend) http.Handler {
	engine := gin.New()

	engine.Use(gin.Recovery())

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			assetlinks := v1.Group("/assetlinks")
			{
				assetlinks.POST("", assetlinksHandler(alb, wellknown.CreateAssetlink))
				assetlinks.PUT("", assetlinksHandler(alb, wellknown.UpdateAssetlink))
				assetlinks.PATCH("", assetlinksHandler(alb, wellknown.PatchAssetlink))
			}

			appleAppSiteAssociation := v1.Group("/apple-app-site-association")
			{
				appleAppSiteAssociation.POST("", appleAppSiteAssociationApplinkDetailHandler(aasab, wellknown.CreateAppleAppSiteAssoicationAppLinksDetail))
				appleAppSiteAssociation.PUT("", appleAppSiteAssociationApplinkDetailHandler(aasab, wellknown.UpdateAppleAppSiteAssoicationAppLinksDetail))
				appleAppSiteAssociation.PATCH("", appleAppSiteAssociationApplinkDetailHandler(aasab, wellknown.PatchAppleAppSiteAssoicationAppLinksDetail))
			}
		}
	}

	return engine
}

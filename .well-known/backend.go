package wellknown

import "context"

type AppleAppSiteAssociationBackend interface {
	GetAppleAppSiteAssociation(context.Context) (*AppleAppSiteAssociation, error)
	PutAppleAppSiteAssociation(context.Context, *AppleAppSiteAssociation) error
}

type AssetlinksBackend interface {
	GetAssetlinks(context.Context) ([]Assetlink, error)
	PutAssetlinks(context.Context, []Assetlink) error
}

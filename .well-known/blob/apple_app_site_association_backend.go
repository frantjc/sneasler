package wellknownblob

import (
	"context"
	"encoding/json"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

type AppleAppSiteAssociationBucket struct {
	*blob.Bucket
	Key string
}

func (b *AppleAppSiteAssociationBucket) GetAppleAppSiteAssociation(ctx context.Context) (*wellknown.AppleAppSiteAssociation, error) {
	r, err := b.Bucket.NewReader(ctx, b.Key, nil)

	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			return &wellknown.AppleAppSiteAssociation{
				AppLinks: &wellknown.AppleAppSiteAssociationAppLinks{
					Apps: []string{},
					Details: []wellknown.AppleAppSiteAssociationAppLinksDetail{},
				},
			}, nil
		}

		return nil, err
	}
	defer r.Close()

	a := &wellknown.AppleAppSiteAssociation{}
	return a, json.NewDecoder(r).Decode(a)
}

func (b *AppleAppSiteAssociationBucket) PutAppleAppSiteAssociation(ctx context.Context, a *wellknown.AppleAppSiteAssociation) error {
	w, err := b.Bucket.NewWriter(ctx, b.Key, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	return json.NewEncoder(w).Encode(a)
}

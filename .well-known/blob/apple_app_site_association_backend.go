package wellknownblob

import (
	"context"
	"encoding/json"
	"sync"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

type AppleAppSiteAssociationBucket struct {
	*blob.Bucket
	Key string
	sync.Mutex
}

func (b *AppleAppSiteAssociationBucket) GetAppleAppSiteAssociation(ctx context.Context) (*wellknown.AppleAppSiteAssociation, error) {
	b.Lock()
	defer b.Unlock()

	r, err := b.Bucket.NewReader(ctx, b.Key, nil)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			return &wellknown.AppleAppSiteAssociation{
				AppLinks: &wellknown.AppleAppSiteAssociationAppLinks{
					AppleAppSiteAssociationApps: wellknown.AppleAppSiteAssociationApps{
						Apps: []string{},
					},
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
	b.Lock()
	defer b.Unlock()

	w, err := b.Bucket.NewWriter(ctx, b.Key, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	return json.NewEncoder(w).Encode(a)
}

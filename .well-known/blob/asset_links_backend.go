package wellknownblob

import (
	"context"
	"encoding/json"

	wellknown "github.com/frantjc/sneasler/.well-known"
	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

type AssetlinksBucket struct {
	*blob.Bucket
	Key string
}

func (b *AssetlinksBucket) GetAssetlinks(ctx context.Context) ([]wellknown.Assetlink, error) {
	r, err := b.Bucket.NewReader(ctx, b.Key, nil)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			return []wellknown.Assetlink{}, nil
		}

		return nil, err
	}
	defer r.Close()

	a := []wellknown.Assetlink{}
	return a, json.NewDecoder(r).Decode(&a)
}

func (b *AssetlinksBucket) PutAssetlinks(ctx context.Context, a []wellknown.Assetlink) error {
	w, err := b.Bucket.NewWriter(ctx, b.Key, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	return json.NewEncoder(w).Encode(a)
}

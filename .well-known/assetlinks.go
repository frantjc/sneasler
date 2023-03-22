package wellknown

import (
	"context"
	"fmt"

	"github.com/frantjc/go-fn"
)

type Assetlink struct {
	Relation []string         `json:"relation"`
	Target   *AssetlinkTarget `json:"target"`
}

type AssetlinkTarget struct {
	Namespace              string   `json:"namespace"`
	PackageName            string   `json:"package_name"`
	SHA256CertFingerprints []string `json:"sha256_cert_fingerprints"`
}

func CreateAssetlink(ctx context.Context, b AssetlinksBackend, l *Assetlink) error {
	a, err := b.GetAssetlinks(ctx)
	if err != nil {
		return err
	}

	for _, e := range a {
		if e.Target.PackageName == l.Target.PackageName {
			return fmt.Errorf("assetlink %s already exists", l.Target.PackageName)
		}
	}

	return b.PutAssetlinks(ctx, append(a, *l))
}

func UpdateAssetlink(ctx context.Context, b AssetlinksBackend, l *Assetlink) error {
	a, err := b.GetAssetlinks(ctx)
	if err != nil {
		return err
	}

	for i, e := range a {
		if e.Target.PackageName == l.Target.PackageName {
			a[i] = *l
			return b.PutAssetlinks(ctx, a)
		}
	}

	return fmt.Errorf("assetlink %s does not exist", l.Target.PackageName)
}

func PatchAssetlink(ctx context.Context, b AssetlinksBackend, l *Assetlink) error {
	a, err := b.GetAssetlinks(ctx)
	if err != nil {
		return err
	}

	for i, e := range a {
		if e.Target.PackageName == l.Target.PackageName {
			a[i] = Assetlink{
				Relation: fn.Unique(append(e.Relation, l.Relation...)),
				Target: &AssetlinkTarget{
					PackageName:            l.Target.PackageName,
					Namespace:              l.Target.Namespace,
					SHA256CertFingerprints: fn.Unique(append(e.Target.SHA256CertFingerprints, l.Target.SHA256CertFingerprints...)),
				},
			}
			return b.PutAssetlinks(ctx, a)
		}
	}

	return fmt.Errorf("assetlink %s does not exist", l.Target.PackageName)
}

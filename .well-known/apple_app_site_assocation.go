package wellknown

import (
	"context"
	"fmt"

	"github.com/frantjc/go-fn"
)

type AppleAppSiteAssociation struct {
	ActivityContinuation *AppleAppSiteAssociationApps     `json:"activitycontinuation,omitempty"`
	WebCredentials       *AppleAppSiteAssociationApps     `json:"webcredentials,omitempty"`
	AppLinks             *AppleAppSiteAssociationAppLinks `json:"applinks,omitempty"`
}

type AppleAppSiteAssociationApps struct {
	Apps []string `json:"apps"`
}

type AppleAppSiteAssociationAppLinks struct {
	AppleAppSiteAssociationApps `json:",inline"`
	Details []AppleAppSiteAssociationAppLinksDetail `json:"details"`
}

type AppleAppSiteAssociationAppLinksDetail struct {
	AppID string   `json:"appID"`
	Paths []string `json:"paths"`
}

func CreateAppleAppSiteAssoicationAppLinksDetail(ctx context.Context, b AppleAppSiteAssociationBackend, d *AppleAppSiteAssociationAppLinksDetail) error {
	a, err := b.GetAppleAppSiteAssociation(ctx)
	if err != nil {
		return err
	}

	for _, e := range a.AppLinks.Details {
		if e.AppID == d.AppID {
			return fmt.Errorf("app %s already exists", d.AppID)
		}
	}

	a.AppLinks.Details = append(a.AppLinks.Details, *d)

	return b.PutAppleAppSiteAssociation(ctx, a)
}

func UpdateAppleAppSiteAssoicationAppLinksDetail(ctx context.Context, b AppleAppSiteAssociationBackend, d *AppleAppSiteAssociationAppLinksDetail) error {
	a, err := b.GetAppleAppSiteAssociation(ctx)
	if err != nil {
		return err
	}

	for i, e := range a.AppLinks.Details {
		if e.AppID == d.AppID {
			a.AppLinks.Details[i] = e
			return b.PutAppleAppSiteAssociation(ctx, a)
		}
	}

	return fmt.Errorf("app %s does not exist", d.AppID)
}

func PatchAppleAppSiteAssoicationAppLinksDetail(ctx context.Context, b AppleAppSiteAssociationBackend, d *AppleAppSiteAssociationAppLinksDetail) error {
	a, err := b.GetAppleAppSiteAssociation(ctx)
	if err != nil {
		return err
	}

	for i, e := range a.AppLinks.Details {
		if e.AppID == d.AppID {
			a.AppLinks.Details[i].Paths = fn.Unique(append(e.Paths, d.Paths...))
			return b.PutAppleAppSiteAssociation(ctx, a)
		}
	}

	return fmt.Errorf("app %s does not exist", d.AppID)
}

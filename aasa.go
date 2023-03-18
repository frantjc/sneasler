package sneasler

type AppleAppSiteAssocation struct {
	Applinks *struct {
		Apps    []string `json:"apps"`
		Details []struct {
			AppID string   `json:"appID"`
			Paths []string `json:"paths"`
		}
	} `json:"applinks"`
}

package sneasler

type AssetLink struct {
	Relation []string `json:"relation"`
	Target   *struct {
		Namespace              string   `json:"namespace"`
		PackageName            string   `json:"package_name"`
		SHA256CertFingerprints []string `json:"sha256_cert_fingerprints"`
	} `json:"target"`
}

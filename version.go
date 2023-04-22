package sneasler

import (
	"runtime/debug"
	"strings"
)

// Semver is the semantic version of sneasler.
// Meant to be be overridden at build time.
var Semver = "0.1.2"

// GetSemver returns the semantic version of sneasler as built from
// Semver and debug build info.
func GetSemver() string {
	version := Semver

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		var (
			revision string
			modified bool
		)
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}

		if revision != "" {
			i := len(revision)
			if i > 7 {
				i = 7
			}

			if !strings.Contains(version, revision[:i]) {
				version += "+" + revision[:i]
			}
		}

		if modified {
			version += "*"
		}
	}

	return version
}

func GetImageRef() string {
	return "ghcr.io/frantjc/sneasler:" + Semver
}

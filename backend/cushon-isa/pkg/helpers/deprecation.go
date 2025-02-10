package helpers

import (
	"strings"
)

// Note: Not used but deprecation solutions considered
var deprecatedVersions = map[string]struct{}{
	"/v1/": {},
}

func isDeprecatedVersion(path string) bool {
	for version := range deprecatedVersions {
		if strings.Contains(path, version) {
			return true
		}
	}
	return false
}

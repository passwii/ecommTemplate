package catalog

import (
	"net/url"
	"strings"
)

// storefrontOrigin is a placeholder used only to resolve site-relative paths,
// mirroring the storefront's own URL safety rule.
const storefrontOrigin = "https://storefront.invalid"

// SafeResourceURL accepts external HTTP(S) URLs or absolute paths that resolve
// within this storefront. It mirrors the storefront's resource URL rule so the
// API never returns a stored link the storefront would reject.
func SafeResourceURL(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}

	if strings.HasPrefix(trimmed, "/") {
		if strings.HasPrefix(trimmed, "//") || strings.Contains(trimmed, "\\") || hasControlCharacter(trimmed) {
			return false
		}
		base, err := url.Parse(storefrontOrigin + "/")
		if err != nil {
			return false
		}
		reference, err := url.Parse(trimmed)
		if err != nil {
			return false
		}
		resolved := base.ResolveReference(reference)
		return resolved.Scheme == "https" && resolved.Host == "storefront.invalid"
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(parsed.Scheme)
	return (scheme == "https" || scheme == "http") && parsed.Hostname() != ""
}

func hasControlCharacter(value string) bool {
	for _, r := range value {
		if r < 0x20 {
			return true
		}
	}
	return false
}

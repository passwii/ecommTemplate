package catalog

import "testing"

func TestSafeResourceURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"site absolute path", "/files/spec.pdf", true},
		{"root path", "/", true},
		{"external https", "https://example.com/a.pdf", true},
		{"external http", "http://example.com/a.pdf", true},
		{"empty", "", false},
		{"whitespace only", "   ", false},
		{"protocol relative", "//evil.example.com/a.pdf", false},
		{"backslash in site path", "/files\\..\\secret", false},
		{"site path traversal stays within origin", "/../../etc/passwd", true},
		{"control character", "/files/a\u0000b.pdf", false},
		{"javascript scheme", "javascript:alert(1)", false},
		{"data scheme", "data:text/html;base64,PHNjcmlwdD4=", false},
		{"scheme without host", "https://", false},
		{"relative without leading slash", "files/a.pdf", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SafeResourceURL(tc.in); got != tc.want {
				t.Errorf("SafeResourceURL(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

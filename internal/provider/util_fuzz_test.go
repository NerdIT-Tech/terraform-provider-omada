package provider

import "testing"

func FuzzFirstNonEmpty(f *testing.F) {
	f.Add("", "fallback")
	f.Add("explicit", "")
	f.Add("", "")
	f.Add("explicit", "fallback")

	f.Fuzz(func(t *testing.T, a, b string) {
		got := firstNonEmpty(a, b)
		switch {
		case a != "" && got != a:
			t.Fatalf("firstNonEmpty(%q, %q) = %q, want %q", a, b, got, a)
		case a == "" && got != b:
			t.Fatalf("firstNonEmpty(%q, %q) = %q, want %q", a, b, got, b)
		}
	})
}

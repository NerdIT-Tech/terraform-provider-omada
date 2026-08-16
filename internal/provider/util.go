package provider

// firstNonEmpty returns the first non-empty string in vals, or "" if all are
// empty. Used to let explicit provider configuration values take precedence
// over environment variable fallbacks.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

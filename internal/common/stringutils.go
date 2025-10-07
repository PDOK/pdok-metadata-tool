package common

// TruncateString truncates a string to maxLen and adds "..." if it was truncated.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-3] + "..."
}

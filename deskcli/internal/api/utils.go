package api

import (
	"strings"
)

// splitLines splits a string into non-empty lines.
func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// splitTab splits a line by tab into at most n parts.
func splitTab(s string, n int) []string {
	return strings.SplitN(s, "\t", n)
}

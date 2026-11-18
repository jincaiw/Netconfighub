// Package diffcap bounds the O(n*m) LCS diff computation so a very large
// configuration cannot amplify memory/CPU usage. The threshold is configurable
// via NCH_LCS_MAX_LINES and defaults to 2000 lines.
package diffcap

import (
	"os"
	"strconv"
	"strings"
)

const DefaultMaxLines = 2000

// MaxLines returns the configured per-side line threshold.
func MaxLines() int {
	value := strings.TrimSpace(os.Getenv("NCH_LCS_MAX_LINES"))
	if value == "" {
		return DefaultMaxLines
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return DefaultMaxLines
	}
	return n
}

// Exceeded reports whether computing a full LCS for these two line slices
// would exceed the configured threshold.
func Exceeded(a, b []string) bool {
	limit := MaxLines()
	return len(a) > limit || len(b) > limit
}

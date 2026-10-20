package img

import (
	"fmt"
	"os"
)

// FileSize returns the size of path, or 0 if missing.
func FileSize(path string) int64 {
	st, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return st.Size()
}

// SumFileSizes totals sizes for existing files (missing paths are skipped).
func SumFileSizes(paths []string) int64 {
	var n int64
	for _, p := range paths {
		n += FileSize(p)
	}
	return n
}

// FormatBytes formats a byte count for display (e.g. "1.2 MB").
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	val := float64(b) / float64(div)
	if val >= 100 || exp == 0 {
		return fmt.Sprintf("%.0f %cB", val, "KMGTPE"[exp])
	}
	return fmt.Sprintf("%.1f %cB", val, "KMGTPE"[exp])
}

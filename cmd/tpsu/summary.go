package main

import (
	"fmt"

	"github.com/Cimmerial/tps-utilities/internal/img"
)

func convertDoneMessage(wrote, unchanged, skipped int, before, after int64) string {
	msg := fmt.Sprintf("done — %d wrote, %d unchanged, %d skipped", wrote, unchanged, skipped)
	if before > 0 || after > 0 {
		msg += fmt.Sprintf(" — %s → %s", img.FormatBytes(before), img.FormatBytes(after))
	}
	return msg
}

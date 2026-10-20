package tui

import (
	"fmt"
	"os"
)

// MenuRow is one line in a pipe-aligned menu (num | description).
type MenuRow struct {
	Num  string
	Desc string
}

// PrintPipeMenu prints rows like TPS help: "    1 | JPEG" with aligned numbers.
func PrintPipeMenu(rows []MenuRow) {
	maxNum := 0
	for _, r := range rows {
		if len(r.Num) > maxNum {
			maxNum = len(r.Num)
		}
	}
	for _, r := range rows {
		numPlain := padPlain(r.Num, maxNum)
		fmt.Fprintf(os.Stderr, "    %s | %s\n", Bold(Abbrev(numPlain)), r.Desc)
	}
}

func padPlain(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + spaces(width-len(s))
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

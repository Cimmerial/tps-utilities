package tui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	red     = "\033[31m"
)

func colorEnabled() bool {
	if colorForcedOff {
		return false
	}
	if v := strings.TrimSpace(os.Getenv("TPSU_COLOR")); v == "1" || strings.EqualFold(v, "true") {
		return true
	}
	if v := strings.TrimSpace(os.Getenv("NO_COLOR")); v != "" {
		return false
	}
	if term.IsTerminal(int(os.Stdout.Fd())) || term.IsTerminal(int(os.Stderr.Fd())) {
		return true
	}
	return false
}

func paint(code, s string) string {
	if !colorEnabled() {
		return s
	}
	return code + s + reset
}

func Bold(s string) string    { return paint(bold, s) }
func Dim(s string) string     { return paint(dim, s) }
func Cyan(s string) string    { return paint(cyan, s) }
func Green(s string) string   { return paint(green, s) }
func Red(s string) string     { return paint(red, s) }
func Yellow(s string) string  { return paint(yellow, s) }
func Magenta(s string) string { return paint(magenta, s) }

func Abbrev(s string) string { return paint(abbrevCode(), s) }

// Header is a top-level section title (bold accent, ALL CAPS) — matches tps help.
func Header(s string) string {
	return Bold(paint(accentCode(), strings.ToUpper(strings.TrimSpace(s))))
}

// SectionLabel is an inline heading (bold accent, keeps casing).
func SectionLabel(label string) string {
	return Bold(paint(accentCode(), strings.TrimSpace(label)))
}

// SubHeader is an indented section label (bold accent, ALL CAPS).
func SubHeader(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		label = "general"
	}
	return Bold(paint(accentCode(), "    "+strings.ToUpper(label)))
}

// Result prints a prominent success (green) or failure (red) line on stderr.
func Result(ok bool, msg string) {
	fmt.Fprintln(os.Stderr)
	line := "FAILED: " + msg
	if ok {
		line = "SUCCESS: " + msg
	}
	fmt.Fprintln(os.Stderr, "    "+resultPaint(ok, line))
	fmt.Fprintln(os.Stderr)
}

func resultPaint(ok bool, s string) string {
	if ok {
		return Green(s)
	}
	return Red(s)
}

func Section(title string) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, Bold(title))
}

func NumberedRow(num string, desc string) {
	fmt.Fprintf(os.Stderr, "    %s  %s\n", Bold(Abbrev(padLeft(num, 1))), desc)
}

func padLeft(s string, min int) string {
	for len(s) < min {
		s = " " + s
	}
	return s
}

func Bullet(name string) {
	fmt.Fprintf(os.Stderr, "    %s\n", Dim(name))
}

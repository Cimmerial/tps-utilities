package tui

import "strings"

var colorPresets = map[string]string{
	"cyan":          "\033[36m",
	"blue":          "\033[34m",
	"magenta":       "\033[35m",
	"red":           "\033[31m",
	"green":         "\033[32m",
	"yellow":        "\033[33m",
	"white":         "\033[37m",
	"bright_cyan":   "\033[96m",
	"bright_blue":   "\033[94m",
	"bright_magenta": "\033[95m",
	"bright_red":    "\033[91m",
	"bright_green":  "\033[92m",
	"bright_yellow": "\033[93m",
	"bright_white":  "\033[97m",
	"orange":        "\033[38;5;208m",
	"purple":        "\033[38;5;135m",
	"pink":          "\033[38;5;213m",
	"lime":          "\033[92m",
	"teal":          "\033[38;5;37m",
}

var colorAlias = map[string]string{
	"default":     "cyan",
	"light_cyan":  "bright_cyan",
	"light_blue":  "bright_blue",
	"light_green": "bright_green",
	"light_red":   "bright_red",
	"gray":        "white",
	"grey":        "white",
}

var abbrevBright = map[string]string{
	"red":     "\033[91m",
	"green":   "\033[92m",
	"yellow":  "\033[93m",
	"blue":    "\033[94m",
	"magenta": "\033[95m",
	"cyan":    "\033[96m",
	"white":   "\033[97m",
}

func normalizeColorID(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "cyan"
	}
	if canon, ok := colorAlias[s]; ok {
		s = canon
	}
	if colorPresets[s] != "" {
		return s
	}
	return "cyan"
}

func accentCode() string {
	return colorPresets[accentID]
}

func abbrevCode() string {
	if code, ok := abbrevBright[accentID]; ok {
		return code
	}
	return accentCode()
}

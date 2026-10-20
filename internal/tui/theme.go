package tui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var accentID = "cyan"
var colorForcedOff bool

type tpsPreferences struct {
	Accent       string `yaml:"accent,omitempty"`
	ColorEnabled *bool  `yaml:"color_enabled,omitempty"`
}

// InitTheme loads accent/color settings from TPS preferences or TPSU_ACCENT.
func InitTheme() {
	accentID = "cyan"
	colorForcedOff = false
	if v := strings.TrimSpace(os.Getenv("TPSU_ACCENT")); v != "" {
		accentID = normalizeColorID(v)
	}
	if prefs, err := loadTPSPreferences(); err == nil {
		if v := strings.TrimSpace(prefs.Accent); v != "" && os.Getenv("TPSU_ACCENT") == "" {
			accentID = normalizeColorID(v)
		}
		if prefs.ColorEnabled != nil && !*prefs.ColorEnabled {
			colorForcedOff = true
		}
	}
}

func loadTPSPreferences() (tpsPreferences, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return tpsPreferences{}, err
	}
	path := filepath.Join(dir, "tps", "preferences.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return tpsPreferences{}, nil
		}
		return tpsPreferences{}, err
	}
	var prefs tpsPreferences
	if err := yaml.Unmarshal(data, &prefs); err != nil {
		return tpsPreferences{}, err
	}
	return prefs, nil
}

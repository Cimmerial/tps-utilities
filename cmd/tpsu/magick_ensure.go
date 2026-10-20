package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Cimmerial/tps-utilities/internal/img"
	"github.com/Cimmerial/tps-utilities/internal/tui"
)

func ensureMagick(in *bufio.Reader) error {
	if img.MagickOnPath() {
		img.WarmMagickFormats()
		return nil
	}
	auto := os.Getenv("TPSU_INSTALL_MAGICK") == "1"
	if !auto {
		fmt.Fprintln(os.Stderr, tui.Yellow("ImageMagick (magick) is not installed or not on PATH."))
		if !stdinIsTTY() {
			return fmt.Errorf("non-interactive session: install ImageMagick manually, set TPSU_MAGICK=/path/to/magick, or set TPSU_INSTALL_MAGICK=1")
		}
		if !askYesNo(in, "Install ImageMagick now? [y/N] ", false) {
			return fmt.Errorf("ImageMagick is required for conversion")
		}
	}
	fmt.Fprintln(os.Stderr, tui.Dim("Installing ImageMagick…"))
	if err := img.InstallMagick(); err != nil {
		return err
	}
	img.ResetMagickCache()
	if !img.MagickOnPath() {
		return fmt.Errorf("install finished but magick was not found — open a new terminal or set TPSU_MAGICK to the magick binary")
	}
	img.WarmMagickFormats()
	tui.Result(true, "ImageMagick ready")
	return nil
}

func stdinIsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

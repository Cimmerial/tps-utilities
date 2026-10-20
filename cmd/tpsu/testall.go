package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps-utilities/internal/img"
	"github.com/Cimmerial/tps-utilities/internal/tui"
)

func runTestall() error {
	fmt.Fprintln(os.Stderr, tui.Dim("testall: go test ./..."))
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repoRoot()
	if err := cmd.Run(); err != nil {
		tui.Result(false, "unit tests failed")
		return err
	}

	if !img.MagickOnPath() {
		tui.Result(true, "unit tests passed (skipped magick integration — magick not on PATH)")
		return nil
	}

	dir, err := os.MkdirTemp("", "tpsu-testall-")
	if err != nil {
		return failErr(err)
	}
	defer os.RemoveAll(dir)

	png := filepath.Join(dir, "hub.png")
	if out, err := exec.Command(img.MagickBin(), "-size", "8x8", "xc:blue", png).CombinedOutput(); err != nil {
		return failErr(fmt.Errorf("create png: %v %s", err, strings.TrimSpace(string(out))))
	}

	var failed []string
	for _, to := range []img.FormatID{img.FormatJPEG, img.FormatPNG, img.FormatGIF, img.FormatWebP} {
		if !img.ProbeWrite(to) {
			continue
		}
		out, err := img.ConvertFile(png, dir, to, img.ConvertOptions{
			IfExists: func(dest, src string) (bool, error) { return true, nil },
		})
		if !convertOK(out, err) {
			failed = append(failed, "png->"+string(to))
			continue
		}
	}

	srcJPG := filepath.Join(dir, "from.jpg")
	if _, err := exec.Command(img.MagickBin(), png, srcJPG).CombinedOutput(); err != nil {
		return failErr(fmt.Errorf("make jpg: %v", err))
	}
	out, err := img.ConvertFile(srcJPG, dir, img.FormatPNG, img.ConvertOptions{
		IfExists: func(dest, src string) (bool, error) { return true, nil },
	})
	if !convertOK(out, err) {
		failed = append(failed, "jpg->png")
	}

	res, err := img.MassConvert(dir, dir, img.FormatPNG, img.ConvertOptions{
		IfExists: func(dest, src string) (bool, error) { return true, nil },
	})
	if err != nil {
		failed = append(failed, "mass")
	} else if res.Attempted() > 0 && res.OK == 0 {
		failed = append(failed, "mass-ok-count")
	}

	for _, id := range []img.FormatID{img.FormatJPEG, img.FormatPNG, img.FormatGIF, img.FormatWebP} {
		if img.ProbeWrite(id) {
			continue
		}
		failed = append(failed, "probe-write-"+string(id))
	}

	if len(failed) > 0 {
		tui.Result(false, "integration: "+strings.Join(failed, ", "))
		return fmt.Errorf("integration checks failed")
	}
	tui.Result(true, "unit tests and magick integration passed")
	return nil
}

func convertOK(out img.ConvertOutcome, err error) bool {
	if err != nil {
		return false
	}
	return out.Status == img.ConvertOK || out.Status == img.ConvertNoChange
}

func repoRoot() string {
	if root := os.Getenv("TPSU_ROOT"); root != "" {
		return root
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "."
		}
		wd = parent
	}
}

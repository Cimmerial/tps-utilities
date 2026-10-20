package img

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConvertFile converts one image to target using a PNG hub (or copy when formats match).
func ConvertFile(src, outDir string, target FormatID, opts ConvertOptions) (ConvertOutcome, error) {
	srcFmt := FormatFromPath(src)
	if srcFmt == "" {
		return ConvertOutcome{}, fmt.Errorf("unknown source type: %s", src)
	}
	if target == FormatRAW {
		return ConvertOutcome{}, fmt.Errorf("RAW is not a write target")
	}
	if outDir == "" {
		outDir = filepath.Dir(src)
	}
	base := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))
	if opts.Basename != "" {
		base = strings.TrimSuffix(opts.Basename, filepath.Ext(opts.Basename))
		base = strings.TrimSpace(base)
		if base == "" {
			return ConvertOutcome{}, fmt.Errorf("invalid basename")
		}
	}
	dest := filepath.Join(outDir, base+TargetExt(target))

	srcClean := filepath.Clean(src)
	destClean := filepath.Clean(dest)

	if srcFmt == target && pathsEqual(srcClean, destClean) {
		return ConvertOutcome{OutPath: src, Status: ConvertNoChange}, nil
	}

	if fileExists(destClean) && !pathsEqual(srcClean, destClean) {
		ok, err := decideOverwrite(opts.IfExists, destClean, srcClean)
		if err != nil {
			return ConvertOutcome{}, err
		}
		if !ok {
			return ConvertOutcome{Status: ConvertSkipped}, nil
		}
	}

	note := ""
	if srcFmt == FormatGIF && srcFmt != target {
		note = "GIF: first frame only"
	}

	if srcFmt == target {
		if err := copyFile(src, dest); err != nil {
			return ConvertOutcome{}, err
		}
		if opts.DeleteOriginal && !pathsEqual(srcClean, destClean) {
			if err := os.Remove(src); err != nil {
				return ConvertOutcome{OutPath: dest, Note: note, Status: ConvertOK},
					fmt.Errorf("wrote %s but could not delete original: %w", dest, err)
			}
		}
		return ConvertOutcome{OutPath: dest, Note: note, Status: ConvertOK}, nil
	}

	pngPath, cleanup, err := toPNGHub(src, srcFmt)
	if err != nil {
		return ConvertOutcome{}, err
	}
	defer cleanup()

	if target == FormatPNG {
		if err := copyFile(pngPath, dest); err != nil {
			return ConvertOutcome{}, err
		}
	} else if err := fromPNGHub(pngPath, dest, target); err != nil {
		return ConvertOutcome{}, err
	}

	if opts.DeleteOriginal && !pathsEqual(srcClean, destClean) {
		if err := os.Remove(src); err != nil {
			return ConvertOutcome{OutPath: dest, Note: note, Status: ConvertOK},
				fmt.Errorf("wrote %s but could not delete original: %w", dest, err)
		}
	}
	return ConvertOutcome{OutPath: dest, Note: note, Status: ConvertOK}, nil
}

func decideOverwrite(d OverwriteDecider, dest, src string) (bool, error) {
	if d == nil {
		return true, nil
	}
	return d(dest, src)
}

func pathsEqual(a, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func toPNGHub(src string, srcFmt FormatID) (pngPath string, cleanup func(), err error) {
	cleanup = func() {}
	if srcFmt == FormatPNG {
		return src, cleanup, nil
	}
	f, err := os.CreateTemp("", "tpsu-*.png")
	if err != nil {
		return "", cleanup, err
	}
	pngPath = f.Name()
	_ = f.Close()
	cleanup = func() { os.Remove(pngPath) }

	args := []string{src, pngPath}
	if srcFmt == FormatGIF {
		args[0] = src + "[0]"
	}
	if out, err := magickCommand(args...).CombinedOutput(); err != nil {
		return "", cleanup, fmt.Errorf("magick to png: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return pngPath, cleanup, nil
}

func fromPNGHub(pngPath, outPath string, target FormatID) error {
	args := []string{pngPath, outPath}
	if out, err := magickCommand(args...).CombinedOutput(); err != nil {
		return fmt.Errorf("magick from png: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, in, 0o644)
}

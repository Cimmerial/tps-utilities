package img

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var resolvedMagick string

// MagickBin returns magick or legacy convert. Honors TPSU_MAGICK if set.
func MagickBin() string {
	if resolvedMagick != "" {
		return resolvedMagick
	}
	if p := strings.TrimSpace(os.Getenv("TPSU_MAGICK")); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			resolvedMagick = p
			return p
		}
	}
	for _, p := range standardMagickPaths() {
		if executableFile(p) {
			resolvedMagick = p
			return p
		}
	}
	for _, name := range []string{"magick", "convert"} {
		if p, err := exec.LookPath(name); err == nil {
			resolvedMagick = p
			return p
		}
	}
	if p := resolveMagickViaShell(); p != "" {
		resolvedMagick = p
		return p
	}
	return ""
}

func standardMagickPaths() []string {
	var out []string
	for _, dir := range []string{"/opt/homebrew/bin", "/usr/local/bin", "/opt/local/bin"} {
		out = append(out, filepath.Join(dir, "magick"), filepath.Join(dir, "convert"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		out = append(out, filepath.Join(home, "bin", "magick"))
	}
	return out
}

func executableFile(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func resolveMagickViaShell() string {
	path := augmentPath(os.Getenv("PATH"))
	for _, name := range []string{"magick", "convert"} {
		cmd := exec.Command("sh", "-c", "command -v "+name)
		cmd.Env = append(os.Environ(), "PATH="+path)
		out, err := cmd.Output()
		if err != nil {
			continue
		}
		p := strings.TrimSpace(string(out))
		if p != "" && executableFile(p) {
			return p
		}
	}
	return ""
}

func augmentPath(path string) string {
	extra := "/opt/homebrew/bin:/usr/local/bin:/opt/local/bin"
	if path == "" {
		return extra + ":/usr/bin:/bin"
	}
	if strings.Contains(path, "/opt/homebrew/bin") {
		return path
	}
	return extra + ":" + path
}

func MagickOnPath() bool {
	return MagickBin() != ""
}

func magickCommand(args ...string) *exec.Cmd {
	return exec.Command(MagickBin(), args...)
}

func magickFormatNames(id FormatID) []string {
	switch id {
	case FormatJPEG:
		return []string{"JPEG", "JPG", "JPE"}
	case FormatTIFF:
		return []string{"TIFF", "TIF"}
	case FormatWebP:
		return []string{"WEBP"}
	case FormatJP2:
		return []string{"JP2", "J2K"}
	case FormatJXL:
		return []string{"JXL"}
	case FormatICO:
		return []string{"ICO", "ICON"}
	case FormatPPM:
		return []string{"PPM"}
	case FormatTGA:
		return []string{"TGA"}
	default:
		return []string{strings.ToUpper(string(id))}
	}
}

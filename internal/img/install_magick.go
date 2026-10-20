package img

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ResetMagickCache clears cached magick binary path (after install).
func ResetMagickCache() {
	resolvedMagick = ""
	resetProbeCaches()
}

// InstallMagick installs ImageMagick using the platform package manager when possible.
func InstallMagick() error {
	switch runtime.GOOS {
	case "darwin":
		return installMagickDarwin()
	case "linux":
		return installMagickLinux()
	case "windows":
		return installMagickWindows()
	default:
		return fmt.Errorf("automatic ImageMagick install is not supported on %s", runtime.GOOS)
	}
}

func installMagickDarwin() error {
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew not found — install from https://brew.sh then run: brew install imagemagick")
	}
	return runInstallCmd(exec.Command("brew", "install", "imagemagick"))
}

func installMagickLinux() error {
	try := []struct {
		bin  string
		args []string
		sudo bool
	}{
		{"apt-get", []string{"install", "-y", "imagemagick"}, true},
		{"dnf", []string{"install", "-y", "ImageMagick"}, true},
		{"yum", []string{"install", "-y", "ImageMagick"}, true},
		{"pacman", []string{"-S", "--noconfirm", "imagemagick"}, true},
		{"zypper", []string{"install", "-y", "ImageMagick"}, true},
		{"apk", []string{"add", "--no-cache", "imagemagick"}, true},
		{"brew", []string{"install", "imagemagick"}, false},
	}
	var last error
	for _, t := range try {
		if _, err := exec.LookPath(t.bin); err != nil {
			continue
		}
		cmd := exec.Command(t.bin, t.args...)
		if t.sudo {
			if err := runInstallCmd(cmd); err == nil {
				return nil
			} else {
				last = err
			}
			sudo := exec.Command("sudo", append([]string{t.bin}, t.args...)...)
			if err := runInstallCmd(sudo); err == nil {
				return nil
			} else {
				last = err
			}
			continue
		}
		if err := runInstallCmd(cmd); err == nil {
			return nil
		} else {
			last = err
		}
	}
	if last != nil {
		return fmt.Errorf("could not install via package manager: %v (try: sudo apt install imagemagick)", last)
	}
	return fmt.Errorf("no supported package manager found (apt, dnf, pacman, brew, …)")
}

func installMagickWindows() error {
	try := []struct {
		name string
		args []string
	}{
		{"winget", []string{"install", "-e", "--id", "ImageMagick.ImageMagick", "--accept-package-agreements", "--accept-source-agreements"}},
		{"choco", []string{"install", "imagemagick", "-y"}},
		{"scoop", []string{"install", "imagemagick"}},
	}
	var last error
	for _, t := range try {
		if _, err := exec.LookPath(t.name); err != nil {
			continue
		}
		if err := runInstallCmd(exec.Command(t.name, t.args...)); err == nil {
			return nil
		} else {
			last = err
		}
	}
	if last != nil {
		return fmt.Errorf("could not install via winget/choco/scoop: %v", last)
	}
	return fmt.Errorf("install winget, Chocolatey, or Scoop, then retry — or install from https://imagemagick.org")
}

func runInstallCmd(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), "PATH="+augmentPath(os.Getenv("PATH")))
	return cmd.Run()
}

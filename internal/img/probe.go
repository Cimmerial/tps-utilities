package img

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	probeCache         sync.Map
	probeCacheGen      int
	formatListOnce     sync.Once
	formatModes        map[string]string
)

// WarmMagickFormats loads ImageMagick format support once and fills the probe cache.
func WarmMagickFormats() {
	if MagickBin() == "" {
		return
	}
	loadMagickFormatList()
	for _, t := range WriteTargets() {
		probeCached(t.ID, false)
	}
}

func ProbeRead(id FormatID) bool {
	return probeCached(id, true)
}

func ProbeWrite(id FormatID) bool {
	if id == FormatRAW || id == FormatSVG {
		return false
	}
	return probeCached(id, false)
}

func probeCached(id FormatID, read bool) bool {
	key := probeKey(id, read)
	if v, ok := probeCache.Load(key); ok {
		return v.(bool)
	}
	ok := probeMagick(id, read)
	probeCache.Store(key, ok)
	return ok
}

func probeKey(id FormatID, read bool) string {
	key := string(id)
	if read {
		key = "r:" + key
	} else {
		key = "w:" + key
	}
	return fmt.Sprintf("%d:%s", probeCacheGen, key)
}

func resetProbeCaches() {
	probeCacheGen++
	formatListOnce = sync.Once{}
	formatModes = nil
}

func loadMagickFormatList() {
	formatListOnce.Do(func() {
		formatModes = make(map[string]string)
		out, err := magickCommand("-list", "format").Output()
		if err != nil {
			return
		}
		for _, line := range strings.Split(string(out), "\n") {
			name, mode, ok := parseMagickFormatLine(line)
			if !ok {
				continue
			}
			key := strings.TrimSuffix(strings.ToUpper(strings.TrimSpace(name)), "*")
			formatModes[key] = mode
		}
	})
}

func probeMagick(id FormatID, read bool) bool {
	bin := MagickBin()
	if bin == "" {
		return false
	}
	if magickRWFromList(id, read) {
		return true
	}
	if read {
		return probeReadSmoke(id)
	}
	return probeWriteSmoke(id)
}

func magickRWFromList(id FormatID, read bool) bool {
	loadMagickFormatList()
	if len(formatModes) == 0 {
		return false
	}
	for _, n := range magickFormatNames(id) {
		mode, ok := formatModes[n]
		if !ok {
			continue
		}
		if read {
			return strings.Contains(mode, "r")
		}
		return strings.Contains(mode, "w")
	}
	return false
}

func formatNameMatches(name string, names []string) bool {
	name = strings.TrimSuffix(strings.ToUpper(strings.TrimSpace(name)), "*")
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}

func parseMagickFormatLine(line string) (formatName, mode string, ok bool) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", "", false
	}
	if strings.EqualFold(fields[0], "Format") || strings.HasPrefix(fields[0], "-") {
		return "", "", false
	}
	formatName = fields[0]
	for i := 1; i < len(fields); i++ {
		if isModeToken(fields[i]) {
			return formatName, fields[i], true
		}
	}
	return "", "", false
}

func isModeToken(s string) bool {
	if len(s) < 2 || len(s) > 4 {
		return false
	}
	for _, c := range s {
		if c != 'r' && c != 'w' && c != '-' && c != '+' {
			return false
		}
	}
	return strings.Contains(s, "r") || strings.Contains(s, "w")
}

func probeWriteSmoke(id FormatID) bool {
	dir, err := os.MkdirTemp("", "tpsu-probe-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	src := filepath.Join(dir, "in.png")
	dst := filepath.Join(dir, "out"+TargetExt(id))
	if _, err := magickCommand("-size", "1x1", "xc:white", src).CombinedOutput(); err != nil {
		return false
	}
	if _, err := magickCommand(src, dst).CombinedOutput(); err != nil {
		return false
	}
	_, err = os.Stat(dst)
	return err == nil
}

func probeReadSmoke(id FormatID) bool {
	dir, err := os.MkdirTemp("", "tpsu-probe-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	ext := TargetExt(id)
	if id == FormatJPEG {
		ext = ".jpg"
	}
	src := filepath.Join(dir, "in"+ext)
	png := filepath.Join(dir, "out.png")
	if _, err := magickCommand("-size", "1x1", "xc:white", src).CombinedOutput(); err != nil {
		return false
	}
	if _, err := magickCommand(src, png).CombinedOutput(); err != nil {
		return false
	}
	_, err = os.Stat(png)
	return err == nil
}

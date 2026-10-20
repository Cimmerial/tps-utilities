package stats

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// DefaultRelativePath is recorded in menus/uic.yaml for documentation.
const DefaultRelativePath = "stats/conversions.yaml"

type file struct {
	Conversions map[string]int `yaml:"conversions"`
}

var mu sync.Mutex

func path() string {
	if p := strings.TrimSpace(os.Getenv("TPSU_STATS_PATH")); p != "" {
		return p
	}
	root := strings.TrimSpace(os.Getenv("TPSU_ROOT"))
	if root == "" {
		if wd, err := os.Getwd(); err == nil {
			root = wd
		}
	}
	return filepath.Join(root, DefaultRelativePath)
}

// Counts returns conversion totals by format id.
func Counts() map[string]int {
	mu.Lock()
	defer mu.Unlock()
	f, err := readFile()
	if err != nil {
		return map[string]int{}
	}
	out := make(map[string]int, len(f.Conversions))
	for k, v := range f.Conversions {
		out[k] = v
	}
	return out
}

// Record adds n successful conversions to format id (e.g. "jpeg").
func Record(format string, n int) error {
	if n <= 0 {
		return nil
	}
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	f, err := readFile()
	if err != nil {
		f = file{Conversions: map[string]int{}}
	}
	if f.Conversions == nil {
		f.Conversions = map[string]int{}
	}
	f.Conversions[format] += n
	return writeFile(f)
}

func readFile() (file, error) {
	p := path()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return file{Conversions: map[string]int{}}, nil
		}
		return file{}, err
	}
	var f file
	if err := yaml.Unmarshal(data, &f); err != nil {
		return file{}, err
	}
	if f.Conversions == nil {
		f.Conversions = map[string]int{}
	}
	return f, nil
}

func writeFile(f file) error {
	p := path()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(&f)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

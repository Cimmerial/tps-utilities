package stats

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecordAndCounts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats.yaml")
	t.Setenv("TPSU_STATS_PATH", path)

	if err := Record("jpeg", 2); err != nil {
		t.Fatal(err)
	}
	if err := Record("png", 1); err != nil {
		t.Fatal(err)
	}
	c := Counts()
	if c["jpeg"] != 2 || c["png"] != 1 {
		t.Fatalf("counts: %v", c)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "jpeg: 2") {
		t.Fatalf("file: %s", data)
	}
}

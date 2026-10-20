package img

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func requireMagick(t *testing.T) {
	if _, err := exec.LookPath("magick"); err != nil {
		t.Skip("magick not on PATH")
	}
}

func testOpts(overwrite bool) ConvertOptions {
	return ConvertOptions{
		IfExists: func(dest, src string) (bool, error) {
			return overwrite, nil
		},
	}
}

func TestHubEachFormat(t *testing.T) {
	requireMagick(t)
	dir := t.TempDir()
	pngPath := filepath.Join(dir, "hub.png")
	if out, err := exec.Command("magick", "-size", "4x4", "xc:red", pngPath).CombinedOutput(); err != nil {
		t.Fatalf("create png: %v %s", err, out)
	}

	targets := []FormatID{FormatJPEG, FormatGIF, FormatWebP, FormatTIFF, FormatBMP, FormatAVIF}
	for _, to := range targets {
		if !ProbeWrite(to) {
			t.Logf("skip write %s", to)
			continue
		}
		out, err := ConvertFile(pngPath, dir, to, testOpts(true))
		if err != nil {
			t.Fatalf("png -> %s: %v", to, err)
		}
		if out.Status != ConvertOK {
			t.Fatalf("png -> %s: status %v", to, out.Status)
		}
		if _, err := os.Stat(out.OutPath); err != nil {
			t.Fatalf("missing output %s: %v", out.OutPath, err)
		}
	}

	sources := []struct {
		ext string
		id  FormatID
	}{
		{".jpg", FormatJPEG},
		{".gif", FormatGIF},
		{".webp", FormatWebP},
	}
	for _, src := range sources {
		if !ProbeRead(src.id) {
			t.Logf("skip read %s", src.id)
			continue
		}
		path := filepath.Join(dir, "src"+src.ext)
		if out, err := exec.Command("magick", pngPath, path).CombinedOutput(); err != nil {
			t.Fatalf("make %s: %v %s", src.ext, err, out)
		}
		out, err := ConvertFile(path, dir, FormatPNG, testOpts(true))
		if err != nil {
			t.Fatalf("%s -> png: %v", src.id, err)
		}
		if _, err := os.Stat(out.OutPath); err != nil {
			t.Fatalf("missing %s", out.OutPath)
		}
	}
}

func TestSameFormatNoChange(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(src, []byte("fake-jpeg"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := ConvertFile(src, dir, FormatJPEG, testOpts(true))
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != ConvertNoChange || out.OutPath != src {
		t.Fatalf("expected unchanged %s, got %+v", src, out)
	}
}

func TestSameFormatCopyRename(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(src, []byte("same-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := ConvertFile(src, dir, FormatJPEG, ConvertOptions{
		Basename: "renamed",
		IfExists: testOpts(true).IfExists,
	})
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "renamed.jpg")
	if out.OutPath != dest {
		t.Fatalf("expected %s, got %s", dest, out.OutPath)
	}
	data, err := os.ReadFile(dest)
	if err != nil || string(data) != "same-bytes" {
		t.Fatalf("copy failed: %v", err)
	}
}

func TestDeclinedOverwriteSkipped(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.jpg")
	existing := filepath.Join(dir, "b.jpg")
	if err := os.WriteFile(src, []byte("src"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existing, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := ConvertFile(src, dir, FormatJPEG, ConvertOptions{
		Basename: "b",
		IfExists: testOpts(false).IfExists,
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != ConvertSkipped {
		t.Fatalf("expected skipped, got %+v", out)
	}
	data, _ := os.ReadFile(existing)
	if string(data) != "old" {
		t.Fatal("dest should be unchanged")
	}
}

func TestMassConvertSummary(t *testing.T) {
	res := MassResult{OK: 8, failureCount: 2, SkippedNonImage: 2, SkippedExists: 1, Failures: []string{"a.jpg", "b.png"}}
	s := FormatMassSummary(res)
	if s == "" || !strings.Contains(s, "8/10") || !strings.Contains(s, "skipped 2 non-images") || !strings.Contains(s, "skipped 1 existing") {
		t.Fatalf("summary: %s", s)
	}
}

func TestMassDeclinedNotFailure(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	if err := os.Mkdir(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "one.jpg")
	existingPNG := filepath.Join(outDir, "one.png")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existingPNG, []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := MassConvert(dir, outDir, FormatPNG, testOpts(false))
	if err != nil {
		t.Fatal(err)
	}
	if res.failureCount != 0 || res.SkippedExists != 1 || res.OK != 0 {
		t.Fatalf("unexpected result %+v", res)
	}
}

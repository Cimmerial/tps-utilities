package img

import "testing"

func TestParseMagickFormatLine(t *testing.T) {
	name, mode, ok := parseMagickFormatLine("     PNG* PNG       rw-   Portable Network Graphics")
	if !ok || name != "PNG*" || mode != "rw-" {
		t.Fatalf("png: name=%q mode=%q ok=%v", name, mode, ok)
	}
	name, mode, ok = parseMagickFormatLine("    JPEG* rw-   Joint Photographic Experts Group")
	if !ok || mode != "rw-" {
		t.Fatalf("jpeg short row: name=%q mode=%q", name, mode)
	}
}

func TestIsModeToken(t *testing.T) {
	if !isModeToken("rw-") {
		t.Fatal("rw-")
	}
	if isModeToken("JPEG") {
		t.Fatal("JPEG is not a mode token")
	}
}

func TestFormatNameMatches(t *testing.T) {
	if !formatNameMatches("PNG*", []string{"PNG"}) {
		t.Fatal("PNG*")
	}
	if formatNameMatches("GIF*", []string{"PNG"}) {
		t.Fatal("gif should not match png")
	}
}

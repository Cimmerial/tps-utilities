package img

import "testing"

func TestFormatBytes(t *testing.T) {
	if FormatBytes(512) != "512 B" {
		t.Fatal(FormatBytes(512))
	}
	if FormatBytes(2048) != "2 KB" {
		t.Fatal(FormatBytes(2048))
	}
}

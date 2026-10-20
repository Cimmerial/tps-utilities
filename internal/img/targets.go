package img

import (
	"github.com/Cimmerial/tps-utilities/internal/stats"
)

// WriteTarget is a format the user may convert into.
type WriteTarget struct {
	ID    FormatID
	Label string
}

// WriteTargets lists conversion targets in menu order (includes formats filtered at runtime).
func WriteTargets() []WriteTarget {
	return []WriteTarget{
		{ID: FormatJPEG, Label: "JPEG"},
		{ID: FormatPNG, Label: "PNG"},
		{ID: FormatGIF, Label: "GIF"},
		{ID: FormatWebP, Label: "WebP"},
		{ID: FormatHEIC, Label: "HEIC"},
		{ID: FormatAVIF, Label: "AVIF"},
		{ID: FormatTIFF, Label: "TIFF"},
		{ID: FormatBMP, Label: "BMP"},
		{ID: FormatICO, Label: "ICO"},
		{ID: FormatJXL, Label: "JPEG XL"},
		{ID: FormatJP2, Label: "JPEG 2000"},
		{ID: FormatTGA, Label: "TGA"},
		{ID: FormatPPM, Label: "PPM"},
	}
}

// WriteTargetRow is one selectable line in the interactive format menu.
type WriteTargetRow struct {
	Target WriteTarget
	Count  int // prior successful conversions (from stats/conversions.yaml)
}

// PrepareWriteTargets returns only formats ImageMagick can write on this system.
func PrepareWriteTargets() []WriteTargetRow {
	WarmMagickFormats()
	counts := stats.Counts()
	var rows []WriteTargetRow
	for _, t := range WriteTargets() {
		if TargetWriteReason(t.ID) != "" {
			continue
		}
		rows = append(rows, WriteTargetRow{
			Target: t,
			Count:  counts[string(t.ID)],
		})
	}
	return rows
}

// TargetWriteReason returns why a target cannot be written, or "" if OK.
func TargetWriteReason(id FormatID) string {
	switch id {
	case FormatRAW:
		return "RAW is not a write target"
	case FormatSVG:
		return "SVG is vector; this tool converts raster images only (use PNG or JPEG)"
	}
	if !ProbeWrite(id) {
		return "ImageMagick cannot write this format on this system"
	}
	return ""
}

package img

import (
	"path/filepath"
	"strings"
)

// FormatID is a lowercase hub id (jpeg, png, …).
type FormatID string

const (
	FormatJPEG FormatID = "jpeg"
	FormatPNG  FormatID = "png"
	FormatGIF  FormatID = "gif"
	FormatWebP FormatID = "webp"
	FormatSVG  FormatID = "svg"
	FormatHEIC FormatID = "heic"
	FormatAVIF FormatID = "avif"
	FormatTIFF FormatID = "tiff"
	FormatBMP  FormatID = "bmp"
	FormatICO  FormatID = "ico"
	FormatJXL  FormatID = "jxl"
	FormatJP2  FormatID = "jp2"
	FormatTGA  FormatID = "tga"
	FormatPPM  FormatID = "ppm"
	FormatRAW  FormatID = "raw"
)

var extToFormat = map[string]FormatID{
	".jpg":  FormatJPEG,
	".jpeg": FormatJPEG,
	".png":  FormatPNG,
	".gif":  FormatGIF,
	".webp": FormatWebP,
	".svg":  FormatSVG,
	".heic": FormatHEIC,
	".avif": FormatAVIF,
	".tif":  FormatTIFF,
	".tiff": FormatTIFF,
	".bmp":  FormatBMP,
	".ico":  FormatICO,
	".jxl":  FormatJXL,
	".jp2":  FormatJP2,
	".j2k":  FormatJP2,
	".tga":  FormatTGA,
	".ppm":  FormatPPM,
	".raw":  FormatRAW,
	".cr2":  FormatRAW,
	".nef":  FormatRAW,
	".arw":  FormatRAW,
}

var knownImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".svg": true, ".heic": true, ".avif": true, ".tif": true, ".tiff": true,
	".bmp": true, ".ico": true, ".jxl": true, ".jp2": true, ".j2k": true,
	".tga": true, ".ppm": true,
	".raw": true, ".cr2": true, ".nef": true, ".arw": true,
}

func FormatFromPath(path string) FormatID {
	ext := strings.ToLower(filepath.Ext(path))
	if f, ok := extToFormat[ext]; ok {
		return f
	}
	return ""
}

func IsKnownImageExt(path string) bool {
	return knownImageExts[strings.ToLower(filepath.Ext(path))]
}

func TargetExt(id FormatID) string {
	switch id {
	case FormatJPEG:
		return ".jpg"
	case FormatTIFF:
		return ".tif"
	default:
		return "." + string(id)
	}
}

func ParseFormatID(s string) (FormatID, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return "", errInvalidFormat
	}
	switch FormatID(s) {
	case FormatJPEG, FormatPNG, FormatGIF, FormatWebP, FormatSVG, FormatHEIC, FormatAVIF, FormatTIFF, FormatBMP,
		FormatICO, FormatJXL, FormatJP2, FormatTGA, FormatPPM, FormatRAW:
		return FormatID(s), nil
	default:
		return "", errInvalidFormat
	}
}

var errInvalidFormat = &formatError{}

type formatError struct{}

func (e *formatError) Error() string { return "unknown format" }

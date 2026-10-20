package img

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FolderSummary describes image files in one directory (non-recursive).
type FolderSummary struct {
	ByFormat   map[FormatID]int
	Names      []string
	Skipped    int
	TotalBytes int64
}

// ImagePathsInFolder returns sorted image paths in dir (non-recursive).
func ImagePathsInFolder(dir string) ([]string, error) {
	sum, err := SummarizeFolder(dir)
	if err != nil {
		return nil, err
	}
	dir = strings.TrimSpace(dir)
	out := make([]string, 0, len(sum.Names))
	for _, name := range sum.Names {
		out = append(out, filepath.Join(dir, name))
	}
	return out, nil
}

func SummarizeFolder(dir string) (FolderSummary, error) {
	dir = strings.TrimSpace(dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return FolderSummary{}, err
	}
	sum := FolderSummary{ByFormat: make(map[FormatID]int)}
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		path := filepath.Join(dir, ent.Name())
		if !IsKnownImageExt(path) {
			sum.Skipped++
			continue
		}
		id := FormatFromPath(path)
		if id == "" {
			sum.Skipped++
			continue
		}
		sum.ByFormat[id]++
		sum.Names = append(sum.Names, ent.Name())
		if info, err := ent.Info(); err == nil {
			sum.TotalBytes += info.Size()
		}
	}
	sort.Strings(sum.Names)
	return sum, nil
}

// FolderTypeCounts returns a short summary like "2 png, 1 jpeg".
func FolderTypeCounts(sum FolderSummary) string {
	if len(sum.Names) == 0 {
		return "no images"
	}
	var parts []string
	for _, id := range sortedFormatIDs(sum.ByFormat) {
		parts = append(parts, fmt.Sprintf("%d %s", sum.ByFormat[id], id))
	}
	return strings.Join(parts, ", ")
}

func FormatFolderSummary(sum FolderSummary) string {
	if len(sum.Names) == 0 {
		if sum.Skipped > 0 {
			return fmt.Sprintf("no images in folder (%d other file(s) skipped)", sum.Skipped)
		}
		return "no images in folder"
	}
	var parts []string
	ids := sortedFormatIDs(sum.ByFormat)
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d %s", sum.ByFormat[id], id))
	}
	s := strings.Join(parts, ", ")
	if sum.Skipped > 0 {
		s += fmt.Sprintf("; skipped %d non-image(s)", sum.Skipped)
	}
	return s
}

// SourceSummary describes explicit input paths.
type SourceSummary struct {
	ByFormat   map[FormatID]int
	Names      []string
	TotalBytes int64
}

func SummarizeSources(paths []string) (SourceSummary, error) {
	sum := SourceSummary{ByFormat: make(map[FormatID]int)}
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id := FormatFromPath(p)
		if id == "" {
			return SourceSummary{}, fmt.Errorf("unknown image type: %s", filepath.Base(p))
		}
		sum.ByFormat[id]++
		sum.Names = append(sum.Names, filepath.Base(p))
		sum.TotalBytes += FileSize(p)
	}
	sort.Strings(sum.Names)
	return sum, nil
}

func TypeCounts(by map[FormatID]int) string {
	if len(by) == 0 {
		return "no images"
	}
	var parts []string
	for _, id := range sortedFormatIDs(by) {
		parts = append(parts, fmt.Sprintf("%d %s", by[id], id))
	}
	return strings.Join(parts, ", ")
}

func sortedFormatIDs(m map[FormatID]int) []FormatID {
	out := make([]FormatID, 0, len(m))
	for id := range m {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

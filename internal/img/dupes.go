package img

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// DuplicateBasenameMode applies when several sources share the same output basename.
type DuplicateBasenameMode int

const (
	DupBasenameDefault   DuplicateBasenameMode = 0
	DupBasenameRename    DuplicateBasenameMode = 1 // name, name2, name3, …
	DupBasenameFirstOnly DuplicateBasenameMode = 2 // convert first only per basename
)

// BatchPlan assigns per-source basenames or skips for duplicate-name batches.
type BatchPlan struct {
	BasenameBySrc map[string]string
	SkipSrc       map[string]bool
}

// HasBasenameCollisions reports whether two or more paths share the same stem name.
func HasBasenameCollisions(paths []string) bool {
	counts := map[string]int{}
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		base := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		counts[base]++
		if counts[base] > 1 {
			return true
		}
	}
	return false
}

// PlanBatch resolves duplicate stems for a conversion batch.
func PlanBatch(paths []string, mode DuplicateBasenameMode) BatchPlan {
	plan := BatchPlan{
		BasenameBySrc: map[string]string{},
		SkipSrc:       map[string]bool{},
	}
	if mode == DupBasenameDefault {
		return plan
	}
	groups := map[string][]string{}
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		base := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		groups[base] = append(groups[base], p)
	}
	for base, members := range groups {
		if len(members) < 2 {
			continue
		}
		sort.Strings(members)
		switch mode {
		case DupBasenameFirstOnly:
			for i := 1; i < len(members); i++ {
				plan.SkipSrc[members[i]] = true
			}
		case DupBasenameRename:
			for i := 1; i < len(members); i++ {
				plan.BasenameBySrc[members[i]] = base + fmt.Sprintf("%d", i+1)
			}
		}
	}
	return plan
}

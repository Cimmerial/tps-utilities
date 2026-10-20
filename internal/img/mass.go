package img

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type MassResult struct {
	OK              int
	Wrote           int // new output files (excludes unchanged)
	Failures        []string // sample of up to 5
	SkippedNonImage int
	SkippedExists    int
	SkippedDuplicate int
	BeforeBytes      int64
	AfterBytes      int64
	failureCount    int
}

func (r MassResult) Attempted() int {
	return r.OK + r.failureCount
}

func MassConvert(dir, outDir string, target FormatID, opts ConvertOptions) (MassResult, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return MassResult{}, fmt.Errorf("folder required")
	}
	if target == FormatRAW {
		return MassResult{}, fmt.Errorf("RAW is not a write target")
	}
	if outDir == "" {
		outDir = dir
	}
	var res MassResult
	var failures []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return res, err
	}
	var paths []string
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		path := filepath.Join(dir, ent.Name())
		if !IsKnownImageExt(path) {
			res.SkippedNonImage++
			continue
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	plan := opts.BatchPlan
	if opts.DuplicateBasename != DupBasenameDefault {
		plan = PlanBatch(paths, opts.DuplicateBasename)
	}
	for _, path := range paths {
		if plan.SkipSrc != nil && plan.SkipSrc[path] {
			res.SkippedDuplicate++
			continue
		}
		res.BeforeBytes += FileSize(path)
		fileOpts := opts
		fileOpts.Basename = ""
		if plan.BasenameBySrc != nil {
			if b := plan.BasenameBySrc[path]; b != "" {
				fileOpts.Basename = b
			}
		}
		outcome, err := ConvertFile(path, outDir, target, fileOpts)
		if err != nil {
			failures = append(failures, filepath.Base(path))
			res.failureCount++
			continue
		}
		switch outcome.Status {
		case ConvertSkipped:
			res.SkippedExists++
		case ConvertOK:
			res.OK++
			res.Wrote++
			res.AfterBytes += FileSize(outcome.OutPath)
		case ConvertNoChange:
			res.OK++
			res.AfterBytes += FileSize(outcome.OutPath)
		}
	}
	res.Failures = sampleFailures(failures, 5)
	return res, nil
}

func sampleFailures(all []string, n int) []string {
	if len(all) == 0 {
		return nil
	}
	if len(all) <= n {
		return all
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := rng.Perm(len(all))
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, all[idx[i]])
	}
	return out
}

func FormatMassSummary(res MassResult) string {
	attempted := res.Attempted()
	rate := "0%"
	if attempted > 0 {
		rate = fmt.Sprintf("%.0f%%", 100*float64(res.OK)/float64(attempted))
	}
	s := fmt.Sprintf("mass convert: %d/%d ok (%s)", res.OK, attempted, rate)
	if res.SkippedNonImage > 0 {
		s += fmt.Sprintf("; skipped %d non-images", res.SkippedNonImage)
	}
	if res.SkippedExists > 0 {
		s += fmt.Sprintf("; skipped %d existing file(s)", res.SkippedExists)
	}
	if res.SkippedDuplicate > 0 {
		s += fmt.Sprintf("; skipped %d duplicate name(s)", res.SkippedDuplicate)
	}
	if len(res.Failures) > 0 {
		s += fmt.Sprintf("; sample failures: %s", strings.Join(res.Failures, ", "))
	}
	return s
}

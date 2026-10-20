package img

import "testing"

func TestPlanBatchRename(t *testing.T) {
	paths := []string{
		"/x/edit-movie-classes.bmp",
		"/x/edit-movie-classes.heic",
		"/x/import-from-letterboxd.bmp",
		"/x/import-from-letterboxd.heic",
	}
	plan := PlanBatch(paths, DupBasenameRename)
	if plan.BasenameBySrc["/x/edit-movie-classes.heic"] != "edit-movie-classes2" {
		t.Fatalf("heic: %q", plan.BasenameBySrc["/x/edit-movie-classes.heic"])
	}
	if plan.SkipSrc["/x/edit-movie-classes.bmp"] {
		t.Fatal("first should not skip")
	}
}

func TestPlanBatchFirstOnly(t *testing.T) {
	paths := []string{"/a/foo.bmp", "/a/foo.heic"}
	plan := PlanBatch(paths, DupBasenameFirstOnly)
	if !plan.SkipSrc["/a/foo.heic"] {
		t.Fatal("second should skip")
	}
}

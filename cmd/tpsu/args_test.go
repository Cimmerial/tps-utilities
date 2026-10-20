package main

import "testing"

func TestParseConvertArgsOutOptional(t *testing.T) {
	target, out, files, _, err := parseConvertArgs([]string{"/tmp/a.jpg", "/tmp/b.jpg", "--out"})
	if err != nil {
		t.Fatal(err)
	}
	if target != "" || out != "" {
		t.Fatalf("target=%q out=%q", target, out)
	}
	if len(files) != 2 {
		t.Fatalf("files: %v", files)
	}
}

func TestParseConvertArgsOutSet(t *testing.T) {
	_, out, files, _, err := parseConvertArgs([]string{"/tmp/a.jpg", "--out", "/tmp/out"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "/tmp/out" || len(files) != 1 {
		t.Fatalf("out=%q files=%v", out, files)
	}
}

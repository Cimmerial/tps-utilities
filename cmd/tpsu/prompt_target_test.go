package main

import "testing"

func TestMenuChoiceIsBack(t *testing.T) {
	for _, line := range []string{"", "0", "q", "Q", "back", "BACK"} {
		if !menuChoiceIsBack(line) {
			t.Fatalf("%q should be back", line)
		}
	}
	if menuChoiceIsBack("1") || menuChoiceIsBack("png") {
		t.Fatal("should not be back")
	}
}

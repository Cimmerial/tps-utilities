//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var b strings.Builder
	b.WriteString("title: Universal Image Conversion\nitems:\n")
	b.WriteString(`  - id: images
    label: select images
    action:
      picks:
        - kind: files
          prompt: Select images
      shell: '{root}/bin/tpsu img convert [PICK1]'
`)
	b.WriteString(`  - id: folder
    label: select folder
    action:
      picks:
        - kind: folder
          prompt: Folder to convert
      shell: '{root}/bin/tpsu img mass [PICK1]'
`)
	if err := os.WriteFile("menus/uic.yaml", []byte(b.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

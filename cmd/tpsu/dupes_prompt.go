package main

import (
	"bufio"

	"github.com/Cimmerial/tps-utilities/internal/img"
)

func resolveDuplicateBasename(in *bufio.Reader, paths []string) img.DuplicateBasenameMode {
	if !img.HasBasenameCollisions(paths) {
		return img.DupBasenameDefault
	}
	if askYesNo(in, "Several files share the same output name — keep all (add 2, 3, … to names)? [Y/n] ", true) {
		return img.DupBasenameRename
	}
	return img.DupBasenameFirstOnly
}

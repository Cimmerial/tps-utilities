package main

import (
	"fmt"
	"strings"
)

func parseConvertArgs(args []string) (target, outDir string, files []string, flags cliFlags, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--to":
			i++
			if i >= len(args) {
				return "", "", nil, flags, fmt.Errorf("--to requires a format")
			}
			target = args[i]
		case "--out":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				outDir = args[i]
			}
		case "--name":
			i++
			if i >= len(args) {
				return "", "", nil, flags, fmt.Errorf("--name requires a value")
			}
			flags.name = args[i]
		case "--delete-original":
			t := true
			flags.deleteOriginal = &t
		case "--keep-original":
			t := false
			flags.deleteOriginal = &t
		case "--overwrite":
			t := true
			flags.overwrite = &t
		case "--no-overwrite":
			t := false
			flags.overwrite = &t
		default:
			files = append(files, args[i])
		}
	}
	return target, outDir, files, flags, nil
}

func parseMassArgs(args []string) (target, outDir, folder string, flags cliFlags, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--to":
			i++
			if i >= len(args) {
				return "", "", "", flags, fmt.Errorf("--to requires a format")
			}
			target = args[i]
		case "--out":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				outDir = args[i]
			}
		case "--delete-original":
			t := true
			flags.deleteOriginal = &t
		case "--keep-original":
			t := false
			flags.deleteOriginal = &t
		case "--overwrite":
			t := true
			flags.overwrite = &t
		case "--no-overwrite":
			t := false
			flags.overwrite = &t
		default:
			folder = args[i]
		}
	}
	return target, outDir, folder, flags, nil
}

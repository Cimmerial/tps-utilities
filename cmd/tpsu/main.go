package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Cimmerial/tps-utilities/internal/img"
	"github.com/Cimmerial/tps-utilities/internal/stats"
	"github.com/Cimmerial/tps-utilities/internal/tui"
)

func main() {
	tui.InitTheme()
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "testall":
		if err := runTestall(); err != nil {
			os.Exit(1)
		}
	case "img":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		if err := runImg(os.Args[2], os.Args[3:]); err != nil {
			failErr(err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: tpsu img convert [--to <fmt>] [--out dir] [--name base] [--delete-original|--keep-original] [--overwrite|--no-overwrite] <files...>")
	fmt.Fprintln(os.Stderr, "       tpsu img mass [--to <fmt>] [--out dir] [--delete-original|--keep-original] [--overwrite|--no-overwrite] <folder>")
	fmt.Fprintln(os.Stderr, "       tpsu img probe --read|--write <fmt>")
	fmt.Fprintln(os.Stderr, "       tpsu testall")
}

func runImg(sub string, args []string) error {
	switch sub {
	case "convert":
		return runConvert(args)
	case "mass":
		return runMass(args)
	case "probe":
		return runProbe(args)
	default:
		return fmt.Errorf("unknown img subcommand %q", sub)
	}
}

type cliFlags struct {
	deleteOriginal *bool
	overwrite      *bool
	name           string
}

func runConvert(args []string) error {
	target, outDir, files, flags, err := parseConvertArgs(args)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no input files")
	}
	clean := make([]string, 0, len(files))
	for _, f := range files {
		f = strings.TrimSpace(f)
		if f != "" {
			clean = append(clean, f)
		}
	}
	files = clean
	if len(files) == 0 {
		return fmt.Errorf("no input files")
	}
	in := stdinReader()
	if err := ensureMagick(in); err != nil {
		return err
	}
	if err := describeSources(files); err != nil {
		return err
	}
	beforeBytes := img.SumFileSizes(files)
	fmtID, err := resolveWriteTarget(in, target)
	if err != nil {
		if errors.Is(err, errMenuBack) {
			return nil
		}
		return err
	}
	if flags.name != "" && len(files) > 1 {
		return fmt.Errorf("--name only allowed with a single input file")
	}
	dupMode := resolveDuplicateBasename(in, files)
	plan := img.PlanBatch(files, dupMode)
	opts, err := buildConvertOptions(len(files), flags)
	if err != nil {
		return err
	}
	opts.DuplicateBasename = dupMode
	var notes []string
	var okN, skipN, sameN, dupSkipN int
	var outPaths []string
	for _, f := range files {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if plan.SkipSrc[f] {
			dupSkipN++
			reportConvertLine(true, "skipped (duplicate name): "+f)
			continue
		}
		fileOpts := opts
		if len(files) != 1 {
			fileOpts.Basename = ""
		}
		if b := plan.BasenameBySrc[f]; b != "" {
			fileOpts.Basename = b
		}
		outcome, err := img.ConvertFile(f, outDir, fmtID, fileOpts)
		if err != nil {
			reportConvertLine(false, fmt.Sprintf("%s: %v", f, err))
			return failErr(err)
		}
		switch outcome.Status {
		case img.ConvertSkipped:
			skipN++
			reportConvertLine(true, "skipped (existing file): "+f)
		case img.ConvertNoChange:
			sameN++
			reportConvertLine(true, "unchanged: "+outcome.OutPath)
			outPaths = append(outPaths, outcome.OutPath)
		case img.ConvertOK:
			okN++
			reportConvertLine(true, "wrote "+outcome.OutPath)
			outPaths = append(outPaths, outcome.OutPath)
		}
		if outcome.Note != "" {
			notes = append(notes, outcome.Note)
		}
	}
	if len(notes) > 0 {
		fmt.Fprintln(os.Stderr, tui.Dim(strings.Join(dedupe(notes), "; ")))
	}
	if okN+sameN == 0 {
		return failErr(fmt.Errorf("no files converted"))
	}
	_ = stats.Record(string(fmtID), okN)
	reportDone(true, convertDoneMessage(okN, sameN, skipN+dupSkipN, beforeBytes, img.SumFileSizes(outPaths)))
	return nil
}

func buildConvertOptions(fileCount int, flags cliFlags) (img.ConvertOptions, error) {
	in := bufio.NewReader(os.Stdin)

	deleteOrig := false
	if flags.deleteOriginal != nil {
		deleteOrig = *flags.deleteOriginal
	} else {
		deleteOrig = askYesNo(in, "Delete original file(s) after successful conversion? [y/N] ", false)
	}

	basename := flags.name
	if fileCount == 1 && basename == "" && flags.name == "" {
		basename = askLine(in, "Output basename (empty = keep current name): ")
	}

	var overwriteMode *bool
	if flags.overwrite != nil {
		overwriteMode = flags.overwrite
	}

	opts := img.ConvertOptions{
		Basename:       basename,
		DeleteOriginal: deleteOrig,
		IfExists:       overwriteDecider(in, overwriteMode),
	}
	return opts, nil
}

func overwriteDecider(in *bufio.Reader, mode *bool) img.OverwriteDecider {
	if mode != nil {
		fixed := *mode
		return func(dest, src string) (bool, error) {
			return fixed, nil
		}
	}
	return func(dest, src string) (bool, error) {
		return askYesNo(in, fmt.Sprintf("Overwrite %s? [y/N] ", dest), false), nil
	}
}

func runMass(args []string) error {
	target, outDir, folder, flags, err := parseMassArgs(args)
	if err != nil {
		return err
	}
	if folder == "" {
		return fmt.Errorf("folder required")
	}
	in := stdinReader()
	if err := ensureMagick(in); err != nil {
		return err
	}
	if err := describeFolder(folder); err != nil {
		return err
	}
	fmtID, err := resolveWriteTarget(in, target)
	if err != nil {
		if errors.Is(err, errMenuBack) {
			return nil
		}
		return err
	}
	paths, err := img.ImagePathsInFolder(folder)
	if err != nil {
		return err
	}
	dupMode := resolveDuplicateBasename(in, paths)
	opts, err := buildConvertOptions(0, flags)
	if err != nil {
		return err
	}
	opts.DuplicateBasename = dupMode
	res, err := img.MassConvert(folder, outDir, fmtID, opts)
	if err != nil {
		return failErr(err)
	}
	summary := img.FormatMassSummary(res)
	fmt.Fprintln(os.Stderr, tui.Dim(summary))
	if res.Attempted() > 0 && res.OK == 0 {
		return failErr(fmt.Errorf("mass convert: 0/%d ok", res.Attempted()))
	}
	_ = stats.Record(string(fmtID), res.Wrote)
	sizeNote := ""
	if res.BeforeBytes > 0 || res.AfterBytes > 0 {
		sizeNote = fmt.Sprintf(" — %s → %s", img.FormatBytes(res.BeforeBytes), img.FormatBytes(res.AfterBytes))
	}
	reportDone(true, summary+sizeNote)
	return nil
}

func askYesNo(in *bufio.Reader, prompt string, defaultNo bool) bool {
	fmt.Fprint(os.Stderr, prompt)
	line, err := in.ReadString('\n')
	if err != nil {
		return defaultNo
	}
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" {
		return defaultNo
	}
	return line == "y" || line == "yes"
}

func askLine(in *bufio.Reader, prompt string) string {
	fmt.Fprint(os.Stderr, prompt)
	line, err := in.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

func dedupe(in []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func runProbe(args []string) error {
	var read, write bool
	var format string
	for _, a := range args {
		switch a {
		case "--read":
			read = true
		case "--write":
			write = true
		default:
			format = a
		}
	}
	id, err := img.ParseFormatID(format)
	if err != nil {
		return err
	}
	if read {
		if !img.ProbeRead(id) {
			os.Exit(1)
		}
		return nil
	}
	if write {
		if !img.ProbeWrite(id) {
			os.Exit(1)
		}
		return nil
	}
	return fmt.Errorf("probe needs --read or --write")
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps-utilities/internal/img"
	"github.com/Cimmerial/tps-utilities/internal/tui"
)

var errMenuBack = errors.New("menu back")

func resolveWriteTarget(in *bufio.Reader, flag string) (img.FormatID, error) {
	flag = strings.TrimSpace(flag)
	if flag != "" {
		id, err := img.ParseFormatID(flag)
		if err != nil {
			return "", err
		}
		if reason := img.TargetWriteReason(id); reason != "" {
			return "", fmt.Errorf("%s", reason)
		}
		return id, nil
	}
	return promptWriteTarget(in)
}

func promptWriteTarget(in *bufio.Reader) (img.FormatID, error) {
	rows := img.PrepareWriteTargets()
	if len(rows) == 0 {
		return "", fmt.Errorf("no output formats available (ImageMagick cannot write any supported types)")
	}
	for {
		printWriteTargetMenu(rows)
		line, err := in.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if menuChoiceIsBack(line) {
			return "", errMenuBack
		}
		n, err := strconv.Atoi(line)
		if err != nil || n < 1 || n > len(rows) {
			fmt.Fprintln(os.Stderr, tui.Yellow("Enter a number from the list, or 0 to go back."))
			continue
		}
		return rows[n-1].Target.ID, nil
	}
}

func menuChoiceIsBack(line string) bool {
	if line == "" {
		return true
	}
	switch strings.ToLower(line) {
	case "0", "q", "quit", "back":
		return true
	default:
		return false
	}
}

func printWriteTargetMenu(rows []img.WriteTargetRow) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "    %s\n", tui.Dim("number to choose — Enter, 0, or q = back"))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, tui.SubHeader("convert to"))
	var menuRows []tui.MenuRow
	for i, row := range rows {
		desc := tui.Bold(row.Target.Label)
		if row.Count > 0 {
			desc += " " + tui.Dim(fmt.Sprintf("(%d)", row.Count))
		}
		menuRows = append(menuRows, tui.MenuRow{Num: fmt.Sprintf("%d", i+1), Desc: desc})
	}
	tui.PrintPipeMenu(menuRows)
	fmt.Fprintln(os.Stderr)
	tui.PrintPipeMenu([]tui.MenuRow{{Num: "0", Desc: tui.Dim("back")}})
	fmt.Fprintln(os.Stderr)
	fmt.Fprint(os.Stderr, "Choice (number): ")
}

func describeSources(paths []string) error {
	sum, err := img.SummarizeSources(paths)
	if err != nil {
		return err
	}
	if len(sum.Names) == 0 {
		return fmt.Errorf("no input files")
	}
	printSourceHeader(sum)
	printImageList(sum.Names)
	return nil
}

func describeFolder(dir string) error {
	sum, err := img.SummarizeFolder(dir)
	if err != nil {
		return err
	}
	if len(sum.Names) == 0 {
		return fmt.Errorf("no images to convert in %s", dir)
	}
	base := filepath.Base(dir)
	counts := img.TypeCounts(sum.ByFormat)
	meta := counts
	if sum.TotalBytes > 0 {
		meta = fmt.Sprintf("%s (%s)", counts, img.FormatBytes(sum.TotalBytes))
	}
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "    %s %s\n", tui.SectionLabel(fmt.Sprintf("Folder %s:", base)), tui.Dim(meta))
	printImageList(sum.Names)
	return nil
}

func printSourceHeader(sum img.SourceSummary) {
	counts := img.TypeCounts(sum.ByFormat)
	meta := counts
	if sum.TotalBytes > 0 {
		meta = fmt.Sprintf("%s (%s)", counts, img.FormatBytes(sum.TotalBytes))
	}
	fmt.Fprintln(os.Stderr)
	label := "Selected:"
	if len(sum.Names) == 1 {
		label = sum.Names[0] + ":"
	}
	fmt.Fprintf(os.Stderr, "    %s %s\n", tui.SectionLabel(label), tui.Dim(meta))
}

func printImageList(names []string) {
	if len(names) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, tui.SubHeader("images"))
	limit := len(names)
	if limit > 20 {
		for i := 0; i < 15; i++ {
			tui.Bullet(names[i])
		}
		fmt.Fprintf(os.Stderr, "    %s\n", tui.Dim(fmt.Sprintf("… and %d more", len(names)-15)))
		return
	}
	for _, n := range names {
		tui.Bullet(n)
	}
}

func stdinReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

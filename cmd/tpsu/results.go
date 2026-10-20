package main

import (
	"fmt"
	"os"

	"github.com/Cimmerial/tps-utilities/internal/tui"
)

func reportConvertLine(ok bool, detail string) {
	if ok {
		fmt.Fprintln(os.Stderr, "    "+tui.Green(detail))
	} else {
		fmt.Fprintln(os.Stderr, "    "+tui.Red(detail))
	}
}

func reportDone(ok bool, msg string) {
	tui.Result(ok, msg)
}

func failErr(err error) error {
	tui.Result(false, err.Error())
	return err
}

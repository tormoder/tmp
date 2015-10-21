package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var scanner = bufio.NewScanner(os.Stdin)

type ui struct {
	w io.Writer
}

func (ui *ui) f(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(ui.w, format, a...)
}

func (ui *ui) ln(a ...interface{}) (int, error) {
	return fmt.Fprintln(ui.w, a...)
}

func (ui *ui) sln(a ...interface{}) (int, error) {
	return fmt.Scanln(a...)
}

func (ui *ui) promptForString(stringName string) string {
	ui.f("Enter %s:\n", stringName)
	scanner.Scan()
	input := scanner.Text()
	if err := scanner.Err(); err != nil {
		ui.f("Error reading %s:\n", stringName, err)
		return ui.promptForString(stringName)
	}
	return input
}

// Package output provides styled terminal output utilities.
package output

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var useColors = hasColorSupport()

func hasColorSupport() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ANSI color codes.
const (
	reset  = "\033[0m"
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	dim    = "\033[2m"
	bold   = "\033[1m"
)

func colorize(color, msg string) string {
	if !useColors {
		return msg
	}
	return color + msg + reset
}

// Success prints a success message to stdout.
func Success(msg string) {
	fmt.Println(colorize(green+bold, "✓ "+msg))
}

// Successf prints a formatted success message to stdout.
func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}

// Error prints an error message to stderr.
func Error(msg string) {
	fmt.Fprintln(os.Stderr, colorize(red+bold, "✗ "+msg))
}

// Errorf prints a formatted error message to stderr.
func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

// Warning prints a warning message to stderr.
func Warning(msg string) {
	fmt.Fprintln(os.Stderr, colorize(yellow+bold, "! "+msg))
}

// Warningf prints a formatted warning message to stderr.
func Warningf(format string, args ...any) {
	Warning(fmt.Sprintf(format, args...))
}

// Info prints an info message to stdout.
func Info(msg string) {
	fmt.Println(colorize(blue+bold, "→ "+msg))
}

// Infof prints a formatted info message to stdout.
func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

// Dim prints a dimmed message to stdout.
func Dim(msg string) {
	fmt.Println(colorize(dim, msg))
}

// Section prints a section header.
func Section(title string) {
	fmt.Println()
	fmt.Println(colorize(bold, title))
}

// Step prints a step message with label alignment.
func Step(label, msg string) {
	fmt.Printf("  %-12s %s\n", colorize(dim, label), msg)
}

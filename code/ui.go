package main

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

func printPending(name string, size int64, measurable bool) {
	if measurable {
		fmt.Printf("%s%s[ ] cleaning %s (%s)%s\n", colorBold, colorCyan, name, formatBytes(size), colorReset)
		return
	}
	fmt.Printf("%s%s[ ] cleaning %s%s\n", colorBold, colorCyan, name, colorReset)
}

func printDone(name string) {
	fmt.Printf("%s%s[x] cleaned %s%s\n\n", colorBold, colorGreen, name, colorReset)
}

func printTotal(total int64, measurable bool) {
	if !measurable {
		return
	}
	fmt.Printf("%s%s[x] cleaned (%s)%s\n", colorBold, colorGreen, formatBytes(total), colorReset)
}

func printSkipped(name string) {
	fmt.Printf("%s%s[ ] skipped %s%s\n\n", colorBold, colorYellow, name, colorReset)
}

func confirm(reader *bufio.Reader, name string) bool {
	fmt.Printf("Proceed cleaning %s? [y/N] ", name)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

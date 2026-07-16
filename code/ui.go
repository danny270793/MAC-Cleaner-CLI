package main

import "fmt"

const (
	colorReset = "\033[0m"
	colorCyan  = "\033[36m"
	colorGreen = "\033[32m"
	colorBold  = "\033[1m"
)

func printPending(name string) {
	fmt.Printf("%s%s[ ] cleaning %s%s\n", colorBold, colorCyan, name, colorReset)
}

func printDone(name string) {
	fmt.Printf("%s%s[x] cleaned %s%s\n\n", colorBold, colorGreen, name, colorReset)
}

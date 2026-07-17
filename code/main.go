package main

import (
	"flag"
	"fmt"

	"danny270793/maccleaner/code/cleaners"
)

func main() {
	flag.Usage = func() {
		output := flag.CommandLine.Output()
		fmt.Fprintln(output, "mac-cleaner - free up disk space by clearing caches you don't need")
		fmt.Fprintln(output, "\nUsage:")
		fmt.Fprintln(output, "  maccleaner [flags]")
		fmt.Fprintln(output, "\nFlags:")
		flag.PrintDefaults()
		fmt.Fprintln(output, "\nExamples:")
		fmt.Fprintln(output, "  maccleaner --all")
		fmt.Fprintln(output, "  maccleaner --docker --gradle")
		fmt.Fprintln(output, "  maccleaner --all --dry-run")
	}

	showVersion := flag.Bool("version", false, "print the version and exit")
	all := flag.Bool("all", false, "run every cleaner")
	docker := flag.Bool("docker", false, "clean docker")
	gradle := flag.Bool("gradle", false, "clean gradle caches")
	libraryCaches := flag.Bool("library-caches", false, "clean ~/Library/Caches")
	pubCache := flag.Bool("pub-cache", false, "clean ~/.pub-cache")
	vscodeExtensions := flag.Bool("vscode-extensions", false, "remove outdated versions of ~/.vscode/extensions")
	xcodeDerivedData := flag.Bool("xcode-derived-data", false, "clean ~/Library/Developer/Xcode/DerivedData")
	coreSimulatorCaches := flag.Bool("core-simulator-caches", false, "clean ~/Library/Developer/CoreSimulator/Caches")
	goModCache := flag.Bool("go-mod-cache", false, "clean the go module cache (go clean -modcache)")
	cargoCache := flag.Bool("cargo-cache", false, "clean ~/.cargo/registry/cache and ~/.cargo/registry/src")
	npmCache := flag.Bool("npm-cache", false, "clean ~/.npm")
	pnpmStore := flag.Bool("pnpm-store", false, "prune the pnpm store (pnpm store prune)")
	dryRun := flag.Bool("dry-run", false, "show what would be cleaned without actually cleaning it")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	var selected []Cleaner
	if *all {
		selected = []Cleaner{
			cleaners.Gradle{},
			cleaners.LibraryCaches{},
			cleaners.PubCache{},
			cleaners.VSCodeExtensions{},
			cleaners.XcodeDerivedData{},
			cleaners.CoreSimulatorCaches{},
			cleaners.GoModCache{},
			cleaners.CargoCache{},
			cleaners.NpmCache{},
			cleaners.PnpmStore{},
			cleaners.Docker{},
		}
	} else {
		if *gradle {
			selected = append(selected, cleaners.Gradle{})
		}
		if *libraryCaches {
			selected = append(selected, cleaners.LibraryCaches{})
		}
		if *pubCache {
			selected = append(selected, cleaners.PubCache{})
		}
		if *vscodeExtensions {
			selected = append(selected, cleaners.VSCodeExtensions{})
		}
		if *xcodeDerivedData {
			selected = append(selected, cleaners.XcodeDerivedData{})
		}
		if *coreSimulatorCaches {
			selected = append(selected, cleaners.CoreSimulatorCaches{})
		}
		if *goModCache {
			selected = append(selected, cleaners.GoModCache{})
		}
		if *cargoCache {
			selected = append(selected, cleaners.CargoCache{})
		}
		if *npmCache {
			selected = append(selected, cleaners.NpmCache{})
		}
		if *pnpmStore {
			selected = append(selected, cleaners.PnpmStore{})
		}
		if *docker {
			selected = append(selected, cleaners.Docker{})
		}
	}

	if len(selected) == 0 {
		fmt.Println("no cleaner selected, pass --all or one of --docker --gradle --library-caches --pub-cache --vscode-extensions --xcode-derived-data --core-simulator-caches --go-mod-cache --cargo-cache --npm-cache --pnpm-store")
		flag.Usage()
		return
	}

	var total int64
	var totalMeasurable bool
	for _, cleaner := range selected {
		size, measurable := cleaner.Size()
		printPending(cleaner.Name(), size, measurable)
		if !*dryRun {
			cleaner.Clean()
		}
		printDone(cleaner.Name())
		if measurable {
			total += size
			totalMeasurable = true
		}
	}
	printTotal(total, totalMeasurable)
}

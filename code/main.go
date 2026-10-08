package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintln(output, "  maccleaner --all --auto-approve")
		fmt.Fprintln(output, "  maccleaner --git-clean --gitpath=/Users/you/Github,/Users/you/Gitlab")
		fmt.Fprintln(output, "  maccleaner --all --dry-run --verbose")
	}

	showVersion := flag.Bool("version", false, "print the version and exit")
	all := flag.Bool("all", false, "run every cleaner")
	docker := flag.Bool("docker", false, "clean docker")
	gitClean := flag.Bool("git-clean", false, "remove git-ignored build/dependency folders (node_modules, build, .gradle, ...) found inside git repos under --gitpath")
	gitPath := flag.String("gitpath", "", "comma-separated root paths to scan for git repos, used with --git-clean or --all")
	gradle := flag.Bool("gradle", false, "clean gradle caches")
	libraryCaches := flag.Bool("library-caches", false, "clean ~/Library/Caches")
	pubCache := flag.Bool("pub-cache", false, "clean ~/.pub-cache")
	vscodeExtensions := flag.Bool("vscode-extensions", false, "remove outdated versions of ~/.vscode/extensions")
	xcodeDerivedData := flag.Bool("xcode-derived-data", false, "clean ~/Library/Developer/Xcode/DerivedData")
	coreSimulatorCaches := flag.Bool("core-simulator-caches", false, "clean ~/Library/Developer/CoreSimulator/Caches")
	goModCache := flag.Bool("go-mod-cache", false, "clean the go module cache (go clean -modcache)")
	cargoCache := flag.Bool("cargo-cache", false, "clean ~/.cargo/registry/cache and ~/.cargo/registry/src")
	npmCache := flag.Bool("npm-cache", false, "clean ~/.npm")
	m2Cache := flag.Bool("m2-cache", false, "clean ~/.m2/repository")
	pnpmStore := flag.Bool("pnpm-store", false, "prune the pnpm store (pnpm store prune)")
	dartServerCache := flag.Bool("dart-server-cache", false, "clean ~/.dartServer")
	vaadinCache := flag.Bool("vaadin-cache", false, "clean ~/.vaadin")
	dryRun := flag.Bool("dry-run", false, "show what would be cleaned without actually cleaning it")
	autoApprove := flag.Bool("auto-approve", false, "skip the confirmation prompt before each cleaner")
	flag.Parse()
	verbose := flag.Bool("verbose", false, "show the commands each cleaner will execute before running it")

	if *showVersion {
		fmt.Println(version)
		return
	}

	gitPaths := parseGitPaths(*gitPath)

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
			cleaners.M2Cache{},
			cleaners.Docker{},
			&cleaners.GitClean{Paths: gitPaths},
			cleaners.DartServerCache{},
			cleaners.VaadinCache{},
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
		if *m2Cache {
			selected = append(selected, cleaners.M2Cache{})
		}
		if *docker {
			selected = append(selected, cleaners.Docker{})
		}
		if *gitClean {
			selected = append(selected, &cleaners.GitClean{Paths: gitPaths})
		}
		if *dartServerCache {
			selected = append(selected, cleaners.DartServerCache{})
		}
		if *vaadinCache {
			selected = append(selected, cleaners.VaadinCache{})
		}
	}

	if len(selected) == 0 {
		fmt.Println("no cleaner selected, pass --all or one of --docker --gradle --library-caches --pub-cache --vscode-extensions --xcode-derived-data --core-simulator-caches --go-mod-cache --cargo-cache --npm-cache --pnpm-store --m2-cache --git-clean --dart-server-cache --vaadin-cache")
		flag.Usage()
		return
	}

	reader := bufio.NewReader(os.Stdin)

	var total int64
	var totalMeasurable bool
	for _, cleaner := range selected {
		size, measurable := cleaner.Size()
		printPending(cleaner.Name(), size, measurable)

		if *verbose {
			printCommands(cleaner.Commands())
		}
		if !*dryRun && !*autoApprove && !confirm(reader, cleaner.Name()) {
			printSkipped(cleaner.Name())
			continue
		}

		cleaned, cleanedMeasurable := size, measurable
		if !*dryRun {
			if actual, ok := cleaner.Clean(); ok {
				cleaned, cleanedMeasurable = actual, true
			}
		}
		printDone(cleaner.Name())
		if cleanedMeasurable {
			total += cleaned
			totalMeasurable = true
		}
	}
	printTotal(total, totalMeasurable)
}

func parseGitPaths(value string) []string {
	var paths []string
	for _, part := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			paths = append(paths, trimmed)
		}
	}

	return paths
}

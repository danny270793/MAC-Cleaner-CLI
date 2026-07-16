package main

import (
	"flag"
	"fmt"

	"danny270793/maccleaner/code/cleaners"
)

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	all := flag.Bool("all", false, "run every cleaner")
	docker := flag.Bool("docker", false, "clean docker")
	gradle := flag.Bool("gradle", false, "clean gradle caches")
	libraryCaches := flag.Bool("library-caches", false, "clean ~/Library/Caches")
	pubCache := flag.Bool("pub-cache", false, "clean ~/.pub-cache")
	vscodeExtensions := flag.Bool("vscode-extensions", false, "remove outdated versions of ~/.vscode/extensions")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	var selected []Cleaner
	if *all {
		selected = []Cleaner{cleaners.Gradle{}, cleaners.LibraryCaches{}, cleaners.PubCache{}, cleaners.VSCodeExtensions{}, cleaners.Docker{}}
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
		if *docker {
			selected = append(selected, cleaners.Docker{})
		}
	}

	if len(selected) == 0 {
		fmt.Println("no cleaner selected, pass --all or one of --docker --gradle --library-caches --pub-cache --vscode-extensions")
		flag.Usage()
		return
	}

	for _, cleaner := range selected {
		size, measurable := cleaner.Size()
		printPending(cleaner.Name(), size, measurable)
		cleaner.Clean()
		printDone(cleaner.Name())
	}
}

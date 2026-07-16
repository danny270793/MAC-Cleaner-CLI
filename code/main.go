package main

import "danny270793/maccleaner/code/cleaners"

func main() {
	all := []Cleaner{
		cleaners.Gradle{},
		cleaners.LibraryCaches{},
		cleaners.PubCache{},
		cleaners.Docker{},
	}

	for _, cleaner := range all {
		size, measurable := cleaner.Size()
		printPending(cleaner.Name(), size, measurable)
		cleaner.Clean()
		printDone(cleaner.Name())
	}
}

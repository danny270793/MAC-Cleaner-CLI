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
		printPending(cleaner.Name())
		cleaner.Clean()
		printDone(cleaner.Name())
	}
}

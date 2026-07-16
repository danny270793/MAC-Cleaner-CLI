package main

func main() {
	cleaners := []Cleaner{
		Gradle{},
		LibraryCaches{},
		PubCache{},
		Docker{},
	}

	for _, cleaner := range cleaners {
		cleaner.Clean()
	}
}

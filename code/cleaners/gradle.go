package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type Gradle struct{}

func (Gradle) Name() string {
	return "gradle"
}

func (Gradle) paths() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return []string{
		filepath.Join(home, ".gradle", "caches"),
		filepath.Join(home, ".gradle", "wrapper", "dists"),
	}, nil
}

func (g Gradle) Size() (int64, bool) {
	paths, err := g.paths()
	if err != nil {
		return 0, false
	}

	return sizeOfPaths(paths...)
}

func (g Gradle) Clean() (int64, bool) {
	paths, err := g.paths()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return 0, false
	}

	for _, path := range paths {
		removeContents(path)
	}

	return 0, false
}

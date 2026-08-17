package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type CargoCache struct{}

func (CargoCache) Name() string {
	return "cargo cache"
}

func (CargoCache) paths() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return []string{
		filepath.Join(home, ".cargo", "registry", "cache"),
		filepath.Join(home, ".cargo", "registry", "src"),
	}, nil
}

func (c CargoCache) Size() (int64, bool) {
	paths, err := c.paths()
	if err != nil {
		return 0, false
	}

	return sizeOfPaths(paths...)
}

func (c CargoCache) Clean() (int64, bool) {
	paths, err := c.paths()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return 0, false
	}

	for _, path := range paths {
		removeContents(path)
	}

	return 0, false
}

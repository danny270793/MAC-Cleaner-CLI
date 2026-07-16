package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type CoreSimulatorCaches struct{}

func (CoreSimulatorCaches) Name() string {
	return "core simulator caches"
}

func (CoreSimulatorCaches) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "Developer", "CoreSimulator", "Caches"), nil
}

func (c CoreSimulatorCaches) Size() (int64, bool) {
	path, err := c.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (c CoreSimulatorCaches) Clean() {
	path, err := c.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(path)
}

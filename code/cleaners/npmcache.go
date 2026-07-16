package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type NpmCache struct{}

func (NpmCache) Name() string {
	return "npm cache"
}

func (NpmCache) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".npm"), nil
}

func (n NpmCache) Size() (int64, bool) {
	path, err := n.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (n NpmCache) Clean() {
	path, err := n.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(path)
}

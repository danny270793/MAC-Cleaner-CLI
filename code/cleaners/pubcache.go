package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type PubCache struct{}

func (PubCache) Name() string {
	return "pub cache"
}

func (PubCache) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".pub-cache"), nil
}

func (p PubCache) Size() (int64, bool) {
	path, err := p.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (p PubCache) Clean() {
	path, err := p.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(path)
}

package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type LibraryCaches struct{}

func (LibraryCaches) Name() string {
	return "library caches"
}

func (LibraryCaches) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "Caches"), nil
}

func (l LibraryCaches) Size() (int64, bool) {
	path, err := l.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (l LibraryCaches) Clean() {
	path, err := l.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(path)
}

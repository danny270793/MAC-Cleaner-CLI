package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type M2Cache struct{}

func (M2Cache) Name() string {
	return "maven local repository"
}

func (M2Cache) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".m2", "repository"), nil
}

func (m M2Cache) Size() (int64, bool) {
	path, err := m.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (m M2Cache) Clean() (int64, bool) {
	path, err := m.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return 0, false
	}

	removeContents(path)
	return 0, false
}

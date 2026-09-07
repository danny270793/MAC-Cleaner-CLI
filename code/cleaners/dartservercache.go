package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type DartServerCache struct{}

func (DartServerCache) Name() string {
	return "dart server cache"
}

func (DartServerCache) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".dartServer"), nil
}

func (d DartServerCache) Size() (int64, bool) {
	path, err := d.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (d DartServerCache) Clean() (int64, bool) {
	path, err := d.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return 0, false
	}

	removeContents(path)
	return 0, false
}

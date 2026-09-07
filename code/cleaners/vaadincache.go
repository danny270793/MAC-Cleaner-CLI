package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type VaadinCache struct{}

func (VaadinCache) Name() string {
	return "vaadin cache"
}

func (VaadinCache) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".vaadin"), nil
}

func (v VaadinCache) Size() (int64, bool) {
	path, err := v.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (v VaadinCache) Clean() (int64, bool) {
	path, err := v.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return 0, false
	}

	removeContents(path)
	return 0, false
}

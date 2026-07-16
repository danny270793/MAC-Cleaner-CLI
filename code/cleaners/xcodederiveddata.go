package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type XcodeDerivedData struct{}

func (XcodeDerivedData) Name() string {
	return "xcode derived data"
}

func (XcodeDerivedData) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "Developer", "Xcode", "DerivedData"), nil
}

func (x XcodeDerivedData) Size() (int64, bool) {
	path, err := x.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (x XcodeDerivedData) Clean() {
	path, err := x.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(path)
}

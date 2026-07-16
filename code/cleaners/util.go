package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

func removeContents(path string) {
	entries, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s not found (nothing to clean)\n", path)
		return
	}
	if err != nil {
		fmt.Printf("failed to read %s: %v\n", path, err)
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			fmt.Printf("failed to remove %s: %v\n", entryPath, err)
		}
	}

	fmt.Printf("cleaned %s\n", path)
}

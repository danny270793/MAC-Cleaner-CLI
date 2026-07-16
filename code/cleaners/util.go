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

func dirSize(path string) (int64, bool) {
	if _, err := os.Stat(path); err != nil {
		return 0, false
	}

	var total int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})

	return total, true
}

func sizeOfPaths(paths ...string) (int64, bool) {
	var total int64
	found := false

	for _, path := range paths {
		if size, ok := dirSize(path); ok {
			total += size
			found = true
		}
	}

	return total, found
}

package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type LibraryCaches struct{}

func (LibraryCaches) Clean() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(filepath.Join(home, "Library", "Caches"))
}

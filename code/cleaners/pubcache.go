package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type PubCache struct{}

func (PubCache) Clean() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(filepath.Join(home, ".pub-cache"))
}

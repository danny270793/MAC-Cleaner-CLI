package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
)

type Gradle struct{}

func (Gradle) Clean() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	removeContents(filepath.Join(home, ".gradle", "caches"))
	removeContents(filepath.Join(home, ".gradle", "wrapper", "dists"))
}

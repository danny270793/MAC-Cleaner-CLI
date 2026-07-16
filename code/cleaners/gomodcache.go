package cleaners

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GoModCache struct{}

func (GoModCache) Name() string {
	return "go module cache"
}

func (GoModCache) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "go", "pkg", "mod"), nil
}

func (g GoModCache) Size() (int64, bool) {
	path, err := g.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (GoModCache) Clean() {
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Println("go not installed (nothing to clean)")
		return
	}

	cmd := exec.Command("go", "clean", "-modcache")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run go clean -modcache:", err)
	}
}

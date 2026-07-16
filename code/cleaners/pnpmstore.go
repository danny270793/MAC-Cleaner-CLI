package cleaners

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type PnpmStore struct{}

func (PnpmStore) Name() string {
	return "pnpm store"
}

func (PnpmStore) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "pnpm", "store"), nil
}

func (p PnpmStore) Size() (int64, bool) {
	path, err := p.path()
	if err != nil {
		return 0, false
	}

	return dirSize(path)
}

func (PnpmStore) Clean() {
	if _, err := exec.LookPath("pnpm"); err != nil {
		fmt.Println("pnpm not installed (nothing to clean)")
		return
	}

	cmd := exec.Command("pnpm", "store", "prune")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run pnpm store prune:", err)
	}
}

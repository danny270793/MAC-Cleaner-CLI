package cleaners

import (
	"path/filepath"
	"testing"
)

func TestPnpmStoreName(t *testing.T) {
	if (PnpmStore{}).Name() != "pnpm store" {
		t.Fatalf("unexpected name %q", (PnpmStore{}).Name())
	}
}

func TestPnpmStoreSize(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, "Library", "pnpm", "store", "v3", "files", "aa", "bb"), 250)

	size, ok := PnpmStore{}.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 250 {
		t.Fatalf("expected size 250, got %d", size)
	}
}

func TestPnpmStoreMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	size, ok := PnpmStore{}.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

func TestPnpmStoreCleanWhenNotInstalled(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	output := captureStdout(t, func() {
		PnpmStore{}.Clean()
	})

	want := "pnpm not installed (nothing to clean)\n"
	if output != want {
		t.Fatalf("expected output %q, got %q", want, output)
	}
}

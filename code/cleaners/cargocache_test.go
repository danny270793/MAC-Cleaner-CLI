package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCargoCacheSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".cargo", "registry", "cache", "crate.crate"), 200)
	writeFile(t, filepath.Join(home, ".cargo", "registry", "src", "crate", "lib.rs"), 300)

	c := CargoCache{}

	size, ok := c.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 500 {
		t.Fatalf("expected size 500, got %d", size)
	}

	c.Clean()

	for _, dir := range []string{
		filepath.Join(home, ".cargo", "registry", "cache"),
		filepath.Join(home, ".cargo", "registry", "src"),
	} {
		if _, err := os.Stat(dir); err != nil {
			t.Fatalf("expected %s to still exist: %v", dir, err)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("failed to read %s: %v", dir, err)
		}
		if len(entries) != 0 {
			t.Fatalf("expected %s to be empty after Clean, found %v", dir, entries)
		}
	}
}

func TestCargoCacheNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	c := CargoCache{}
	if c.Name() != "cargo cache" {
		t.Fatalf("unexpected name %q", c.Name())
	}

	size, ok := c.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

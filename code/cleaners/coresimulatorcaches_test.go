package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCoreSimulatorCachesSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, "Library", "Developer", "CoreSimulator", "Caches")
	writeFile(t, filepath.Join(target, "some-cache", "file.bin"), 300)

	c := CoreSimulatorCaches{}

	size, ok := c.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 300 {
		t.Fatalf("expected size 300, got %d", size)
	}

	c.Clean()

	if _, err := os.Stat(target); err != nil {
		t.Fatalf("expected %s to still exist: %v", target, err)
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		t.Fatalf("failed to read %s: %v", target, err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected %s to be empty after Clean, found %v", target, entries)
	}
}

func TestCoreSimulatorCachesNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	c := CoreSimulatorCaches{}
	if c.Name() != "core simulator caches" {
		t.Fatalf("unexpected name %q", c.Name())
	}

	size, ok := c.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

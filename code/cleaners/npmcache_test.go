package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNpmCacheSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, ".npm")
	writeFile(t, filepath.Join(target, "_cacache", "file.bin"), 600)

	n := NpmCache{}

	size, ok := n.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 600 {
		t.Fatalf("expected size 600, got %d", size)
	}

	n.Clean()

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

func TestNpmCacheNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	n := NpmCache{}
	if n.Name() != "npm cache" {
		t.Fatalf("unexpected name %q", n.Name())
	}

	size, ok := n.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

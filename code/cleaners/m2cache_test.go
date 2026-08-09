package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestM2CacheSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, ".m2", "repository")
	writeFile(t, filepath.Join(target, "com", "example", "artifact.jar"), 700)

	m := M2Cache{}

	size, ok := m.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 700 {
		t.Fatalf("expected size 700, got %d", size)
	}

	m.Clean()

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

func TestM2CacheNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	m := M2Cache{}
	if m.Name() != "maven local repository" {
		t.Fatalf("unexpected name %q", m.Name())
	}

	size, ok := m.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

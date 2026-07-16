package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLibraryCachesSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, "Library", "Caches")
	writeFile(t, filepath.Join(target, "com.example.app", "file.bin"), 400)

	lc := LibraryCaches{}

	size, ok := lc.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 400 {
		t.Fatalf("expected size 400, got %d", size)
	}

	lc.Clean()

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

func TestLibraryCachesNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	lc := LibraryCaches{}
	if lc.Name() != "library caches" {
		t.Fatalf("unexpected name %q", lc.Name())
	}

	size, ok := lc.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

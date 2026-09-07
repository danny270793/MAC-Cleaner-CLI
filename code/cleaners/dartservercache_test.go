package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDartServerCacheSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, ".dartServer")
	writeFile(t, filepath.Join(target, "cache", "file.bin"), 400)

	d := DartServerCache{}

	size, ok := d.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 400 {
		t.Fatalf("expected size 400, got %d", size)
	}

	d.Clean()

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

func TestDartServerCacheNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	d := DartServerCache{}
	if d.Name() != "dart server cache" {
		t.Fatalf("unexpected name %q", d.Name())
	}

	size, ok := d.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPubCacheSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, ".pub-cache")
	writeFile(t, filepath.Join(target, "hosted", "file.bin"), 700)

	pc := PubCache{}

	size, ok := pc.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 700 {
		t.Fatalf("expected size 700, got %d", size)
	}

	pc.Clean()

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

func TestPubCacheNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	pc := PubCache{}
	if pc.Name() != "pub cache" {
		t.Fatalf("unexpected name %q", pc.Name())
	}

	size, ok := pc.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

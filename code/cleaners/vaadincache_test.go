package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVaadinCacheSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, ".vaadin")
	writeFile(t, filepath.Join(target, "feature-flags.json"), 300)

	v := VaadinCache{}

	size, ok := v.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 300 {
		t.Fatalf("expected size 300, got %d", size)
	}

	v.Clean()

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

func TestVaadinCacheNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	v := VaadinCache{}
	if v.Name() != "vaadin cache" {
		t.Fatalf("unexpected name %q", v.Name())
	}

	size, ok := v.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

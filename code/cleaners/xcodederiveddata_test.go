package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestXcodeDerivedDataSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target := filepath.Join(home, "Library", "Developer", "Xcode", "DerivedData")
	writeFile(t, filepath.Join(target, "MyApp-abc123", "file.o"), 400)

	x := XcodeDerivedData{}

	size, ok := x.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 400 {
		t.Fatalf("expected size 400, got %d", size)
	}

	x.Clean()

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

func TestXcodeDerivedDataNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	x := XcodeDerivedData{}
	if x.Name() != "xcode derived data" {
		t.Fatalf("unexpected name %q", x.Name())
	}

	size, ok := x.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

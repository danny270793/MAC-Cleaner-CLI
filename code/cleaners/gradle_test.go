package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGradleSizeAndClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".gradle", "caches", "modules", "a.jar"), 200)
	writeFile(t, filepath.Join(home, ".gradle", "wrapper", "dists", "dist", "b.zip"), 300)

	gradle := Gradle{}

	size, ok := gradle.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 500 {
		t.Fatalf("expected size 500, got %d", size)
	}

	gradle.Clean()

	for _, dir := range []string{
		filepath.Join(home, ".gradle", "caches"),
		filepath.Join(home, ".gradle", "wrapper", "dists"),
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

func TestGradleNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	gradle := Gradle{}
	if gradle.Name() != "gradle" {
		t.Fatalf("unexpected name %q", gradle.Name())
	}

	size, ok := gradle.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

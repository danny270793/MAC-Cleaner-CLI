package cleaners

import (
	"path/filepath"
	"testing"
)

func TestGoModCacheName(t *testing.T) {
	if (GoModCache{}).Name() != "go module cache" {
		t.Fatalf("unexpected name %q", (GoModCache{}).Name())
	}
}

func TestGoModCacheSize(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, "go", "pkg", "mod", "example.com", "pkg@v1.0.0", "file.go"), 500)

	size, ok := GoModCache{}.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 500 {
		t.Fatalf("expected size 500, got %d", size)
	}
}

func TestGoModCacheMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	size, ok := GoModCache{}.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

func TestGoModCacheCleanWhenNotInstalled(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	output := captureStdout(t, func() {
		GoModCache{}.Clean()
	})

	want := "go not installed (nothing to clean)\n"
	if output != want {
		t.Fatalf("expected output %q, got %q", want, output)
	}
}

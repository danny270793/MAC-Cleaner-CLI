package cleaners

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path string, size int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	original := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = original }()

	f()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read captured output: %v", err)
	}

	return string(data)
}

func TestDirSize(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), 100)
	writeFile(t, filepath.Join(root, "sub", "b.txt"), 250)

	size, ok := dirSize(root)
	if !ok {
		t.Fatalf("expected dirSize to report ok=true")
	}
	if size != 350 {
		t.Fatalf("expected size 350, got %d", size)
	}
}

func TestDirSizeMissing(t *testing.T) {
	_, ok := dirSize(filepath.Join(t.TempDir(), "missing"))
	if ok {
		t.Fatalf("expected ok=false for a missing directory")
	}
}

func TestSizeOfPaths(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a", "f"), 100)
	writeFile(t, filepath.Join(root, "b", "f"), 50)

	size, ok := sizeOfPaths(filepath.Join(root, "a"), filepath.Join(root, "b"), filepath.Join(root, "missing"))
	if !ok {
		t.Fatalf("expected ok=true when at least one path exists")
	}
	if size != 150 {
		t.Fatalf("expected size 150, got %d", size)
	}
}

func TestSizeOfPathsAllMissing(t *testing.T) {
	root := t.TempDir()
	_, ok := sizeOfPaths(filepath.Join(root, "missing1"), filepath.Join(root, "missing2"))
	if ok {
		t.Fatalf("expected ok=false when no paths exist")
	}
}

func TestRemoveContentsKeepsFolder(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	writeFile(t, filepath.Join(target, "file.txt"), 10)
	writeFile(t, filepath.Join(target, "sub", "nested.txt"), 10)

	removeContents(target)

	if _, err := os.Stat(target); err != nil {
		t.Fatalf("expected folder %s to still exist: %v", target, err)
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		t.Fatalf("failed to read %s: %v", target, err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected %s to be empty, found %v", target, entries)
	}
}

func TestRemoveContentsMissingDoesNotPanic(t *testing.T) {
	removeContents(filepath.Join(t.TempDir(), "missing"))
}

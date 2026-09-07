package cleaners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitCleanFindsAndRemovesJunkFolders(t *testing.T) {
	root := t.TempDir()

	repoA := filepath.Join(root, "org", "repo-a")
	writeFile(t, filepath.Join(repoA, ".git", "HEAD"), 10)
	writeFile(t, filepath.Join(repoA, "node_modules", "left-pad", "index.js"), 100)
	writeFile(t, filepath.Join(repoA, "src", "main.go"), 20)

	repoB := filepath.Join(root, "org", "repo-b")
	writeFile(t, filepath.Join(repoB, ".git", "HEAD"), 10)
	writeFile(t, filepath.Join(repoB, "build", "output.apk"), 200)
	writeFile(t, filepath.Join(repoB, "app", ".gradle", "cache.bin"), 50)

	notARepo := filepath.Join(root, "not-a-repo")
	writeFile(t, filepath.Join(notARepo, "node_modules", "left-pad", "index.js"), 999)

	gitClean := GitClean{Paths: []string{root}}

	size, ok := gitClean.Size()
	if !ok {
		t.Fatalf("expected measurable size")
	}
	if size != 350 {
		t.Fatalf("expected size 350, got %d", size)
	}

	cleaned, ok := gitClean.Clean()
	if !ok {
		t.Fatalf("expected Clean to report ok=true")
	}
	if cleaned != 350 {
		t.Fatalf("expected cleaned 350, got %d", cleaned)
	}

	for _, removed := range []string{
		filepath.Join(repoA, "node_modules"),
		filepath.Join(repoB, "build"),
		filepath.Join(repoB, "app", ".gradle"),
	} {
		if _, err := os.Stat(removed); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, err=%v", removed, err)
		}
	}

	if _, err := os.Stat(filepath.Join(repoA, "src", "main.go")); err != nil {
		t.Fatalf("expected non-junk file to survive: %v", err)
	}
	if _, err := os.Stat(filepath.Join(notARepo, "node_modules")); err != nil {
		t.Fatalf("expected node_modules outside a git repo to survive: %v", err)
	}
}

func TestGitCleanNameAndNoPaths(t *testing.T) {
	gitClean := GitClean{}
	if gitClean.Name() != "git clean" {
		t.Fatalf("unexpected name %q", gitClean.Name())
	}

	size, ok := gitClean.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}

	cleaned, ok := gitClean.Clean()
	if ok {
		t.Fatalf("expected ok=false, got cleaned=%d", cleaned)
	}
}

func TestGitCleanNoMatches(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	writeFile(t, filepath.Join(repo, ".git", "HEAD"), 10)
	writeFile(t, filepath.Join(repo, "src", "main.go"), 20)

	gitClean := GitClean{Paths: []string{root}}

	size, ok := gitClean.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}

	cleaned, ok := gitClean.Clean()
	if ok {
		t.Fatalf("expected ok=false, got cleaned=%d", cleaned)
	}
}

package cleaners

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initRepo creates a real git repo at dir with the given .gitignore content.
func initRepo(t *testing.T, dir, gitignore string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create %s: %v", dir, err)
	}
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v: %s", err, out)
	}
	if gitignore != "" {
		if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(gitignore), 0o644); err != nil {
			t.Fatalf("failed to write .gitignore: %v", err)
		}
	}
}

func TestGitCleanFindsAndRemovesJunkFolders(t *testing.T) {
	root := t.TempDir()

	repoA := filepath.Join(root, "org", "repo-a")
	initRepo(t, repoA, "node_modules/\n")
	writeFile(t, filepath.Join(repoA, "node_modules", "left-pad", "index.js"), 100)
	writeFile(t, filepath.Join(repoA, "src", "main.go"), 20)

	repoB := filepath.Join(root, "org", "repo-b")
	initRepo(t, repoB, "build/\n.gradle/\n")
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
	initRepo(t, repo, "")
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

func TestGitCleanSkipsFoldersNotIgnored(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	initRepo(t, repo, "node_modules/\n")
	writeFile(t, filepath.Join(repo, "node_modules", "pkg", "index.js"), 100)
	writeFile(t, filepath.Join(repo, "dist", "bundle.js"), 200)
	writeFile(t, filepath.Join(repo, "dist", "nested", "node_modules", "x.js"), 30)

	gitClean := GitClean{Paths: []string{root}}

	// dist is not ignored so it must stay, but the ignored node_modules
	// nested inside it is still regenerable.
	size, ok := gitClean.Size()
	if !ok || size != 130 {
		t.Fatalf("expected size 130, got %d (ok=%v)", size, ok)
	}

	if _, ok := gitClean.Clean(); !ok {
		t.Fatalf("expected Clean to report ok=true")
	}

	if _, err := os.Stat(filepath.Join(repo, "dist", "bundle.js")); err != nil {
		t.Fatalf("expected non-ignored dist to survive: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "node_modules")); !os.IsNotExist(err) {
		t.Fatalf("expected ignored node_modules to be removed, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "dist", "nested", "node_modules")); !os.IsNotExist(err) {
		t.Fatalf("expected nested ignored node_modules to be removed, err=%v", err)
	}
}

func TestGitCleanSkipsFoldersWithTrackedFiles(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	initRepo(t, repo, "dist/\n")
	writeFile(t, filepath.Join(repo, "dist", "committed.js"), 100)

	// Force-add a file inside an ignored folder so it is tracked.
	if out, err := exec.Command("git", "-C", repo, "add", "-f", "dist/committed.js").CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v: %s", err, out)
	}

	gitClean := GitClean{Paths: []string{root}}

	if size, ok := gitClean.Size(); ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
	if _, ok := gitClean.Clean(); ok {
		t.Fatalf("expected Clean to report ok=false")
	}
	if _, err := os.Stat(filepath.Join(repo, "dist", "committed.js")); err != nil {
		t.Fatalf("expected tracked file to survive: %v", err)
	}
}

func TestGitCleanCommands(t *testing.T) {
	if cmds := (&GitClean{}).Commands(); len(cmds) != 0 {
		t.Fatalf("expected no commands without paths, got %v", cmds)
	}

	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	initRepo(t, repo, "node_modules/\n")
	writeFile(t, filepath.Join(repo, "node_modules", "pkg", "index.js"), 100)
	writeFile(t, filepath.Join(repo, "dist", "bundle.js"), 200)

	gitClean := GitClean{Paths: []string{root}}
	cmds := gitClean.Commands()
	want := "rm -rf " + shellQuote(filepath.Join(repo, "node_modules"))
	if len(cmds) != 1 || cmds[0] != want {
		t.Fatalf("expected [%q], got %v", want, cmds)
	}
}

func TestGitCleanFindsTargetsInsideIgnoredFolders(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	initRepo(t, repo, "out/\n*.log\n")
	writeFile(t, filepath.Join(repo, "out", "web", "node_modules", "x.js"), 40)
	writeFile(t, filepath.Join(repo, "out", "keep.txt"), 5)

	gitClean := GitClean{Paths: []string{root}}

	cmds := gitClean.Commands()
	want := "rm -rf " + shellQuote(filepath.Join(repo, "out", "web", "node_modules"))
	if len(cmds) != 1 || cmds[0] != want {
		t.Fatalf("expected [%q], got %v", want, cmds)
	}
}

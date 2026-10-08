package cleaners

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// GitClean scans one or more root directories for git repositories and
// removes common build/dependency folders found inside each one, but only
// those that git itself reports as ignored (via .gitignore, .git/info/exclude
// or the global excludes file) and that contain no tracked files.
type GitClean struct {
	Paths []string

	scanned bool
	junk    []string
}

func (GitClean) Name() string {
	return "git clean"
}

// gitCleanTargets are directory names that hold regenerable build or
// dependency output across common ecosystems (Node.js, Java/Gradle/Android,
// Flutter/Dart, Rust/Maven, Next.js).
var gitCleanTargets = []string{
	"node_modules",
	"build",
	".gradle",
	".dart_tool",
	"target",
	".next",
	"dist",
}

func isGitCleanTarget(name string) bool {
	for _, target := range gitCleanTargets {
		if name == target {
			return true
		}
	}
	return false
}

// findRepos walks each root looking for ".git" directories, returning the
// repo root (the parent of ".git") for each one found.
func (g GitClean) findRepos() []string {
	var repos []string
	for _, root := range g.Paths {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}

		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			if d.Name() == ".git" {
				repos = append(repos, filepath.Dir(path))
				return filepath.SkipDir
			}
			return nil
		})
	}

	return repos
}

// ignoredDirs asks git, in a single call, for the directories under repo that
// are ignored and contain no tracked files (a directory holding a tracked
// file is reported file by file instead). Any git failure yields no dirs.
func ignoredDirs(repo string) []string {
	output, err := exec.Command("git", "-C", repo, "ls-files", "-z", "--others", "--ignored", "--exclude-standard", "--directory").Output()
	if err != nil {
		return nil
	}

	var dirs []string
	for _, entry := range strings.Split(string(output), "\x00") {
		if strings.HasSuffix(entry, "/") {
			dirs = append(dirs, filepath.Join(repo, filepath.FromSlash(strings.TrimSuffix(entry, "/"))))
		}
	}

	return dirs
}

// junkPaths returns the paths matching gitCleanTargets that git ignores and
// that hold no tracked files. It scans once and reuses the result, since
// Size, Commands and Clean all need it and git is slow to start.
func (g *GitClean) junkPaths() []string {
	if g.scanned {
		return g.junk
	}

	g.scanned = true
	g.junk = nil
	if _, err := exec.LookPath("git"); err != nil {
		return g.junk
	}

	seen := make(map[string]bool)
	for _, repo := range g.findRepos() {
		for _, dir := range ignoredDirs(repo) {
			for _, target := range targetsWithin(dir) {
				if !seen[target] {
					seen[target] = true
					g.junk = append(g.junk, target)
				}
			}
		}
	}
	g.junk = dropNested(g.junk)

	return g.junk
}

// dropNested removes paths that sit inside another path of the list, so a
// folder is never counted or deleted twice.
func dropNested(paths []string) []string {
	sort.Strings(paths)

	var kept []string
	for _, path := range paths {
		if len(kept) > 0 && strings.HasPrefix(path, kept[len(kept)-1]+string(filepath.Separator)) {
			continue
		}
		kept = append(kept, path)
	}

	return kept
}

// targetsWithin returns dir itself if it is a gitCleanTargets match, else the
// matches nested inside it. dir is wholly ignored with nothing tracked, so
// everything under it is safe too; git reports it as one entry, hiding e.g.
// an "out/node_modules".
func targetsWithin(dir string) []string {
	if isGitCleanTarget(filepath.Base(dir)) {
		return []string{dir}
	}

	var found []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == dir {
			return nil
		}
		if isGitCleanTarget(d.Name()) {
			found = append(found, path)
			return filepath.SkipDir
		}
		return nil
	})

	return found
}

func (g *GitClean) Size() (int64, bool) {
	if len(g.Paths) == 0 {
		return 0, false
	}

	return sizeOfPaths(g.junkPaths()...)
}

func (g *GitClean) Clean() (int64, bool) {
	if len(g.Paths) == 0 {
		fmt.Println("no --gitpath provided (nothing to clean)")
		return 0, false
	}

	paths := g.junkPaths()
	if len(paths) == 0 {
		fmt.Println("no git-ignored build folders found under the given --gitpath (nothing to clean)")
		return 0, false
	}

	var total int64
	for _, path := range paths {
		size, _ := dirSize(path)
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("failed to remove %s: %v\n", path, err)
			continue
		}
		fmt.Printf("cleaned %s\n", path)
		total += size
	}

	// What was scanned is gone now, so a later call must look again.
	g.scanned = false

	return total, true
}

// Commands lists the removals Clean would perform. The read-only git checks
// used to decide what is ignored (git check-ignore, git ls-files) are not
// listed.
func (g *GitClean) Commands() []string {
	if len(g.Paths) == 0 {
		return nil
	}

	var commands []string
	for _, path := range g.junkPaths() {
		commands = append(commands, removeCommand(path))
	}

	return commands
}

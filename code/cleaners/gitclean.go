package cleaners

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// GitClean scans one or more root directories for git repositories and
// removes common build/dependency folders found inside each one.
type GitClean struct {
	Paths []string
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

// junkPaths walks each discovered repo and collects the paths matching
// gitCleanTargets, without descending into a match or into ".git".
func (g GitClean) junkPaths() []string {
	var junk []string
	for _, repo := range g.findRepos() {
		filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			if path != repo && d.Name() == ".git" {
				return filepath.SkipDir
			}
			if isGitCleanTarget(d.Name()) {
				junk = append(junk, path)
				return filepath.SkipDir
			}
			return nil
		})
	}

	return junk
}

func (g GitClean) Size() (int64, bool) {
	if len(g.Paths) == 0 {
		return 0, false
	}

	return sizeOfPaths(g.junkPaths()...)
}

func (g GitClean) Clean() (int64, bool) {
	if len(g.Paths) == 0 {
		fmt.Println("no --gitpath provided (nothing to clean)")
		return 0, false
	}

	paths := g.junkPaths()
	if len(paths) == 0 {
		fmt.Println("no matching build folders found under the given --gitpath (nothing to clean)")
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

	return total, true
}

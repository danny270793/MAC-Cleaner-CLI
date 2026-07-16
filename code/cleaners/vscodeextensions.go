package cleaners

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var vscodeExtensionDirPattern = regexp.MustCompile(`^(.+)-(\d+(?:\.\d+){1,3})(-(?:win32|linux|darwin|alpine|web)(?:-[a-z0-9]+)?)?$`)

type VSCodeExtensions struct{}

func (VSCodeExtensions) Name() string {
	return "vscode extensions"
}

func (VSCodeExtensions) path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".vscode", "extensions"), nil
}

type extensionVersion struct {
	dirName string
	version []int
}

func parseVersion(version string) []int {
	parts := strings.Split(version, ".")
	parsed := make([]int, len(parts))
	for i, part := range parts {
		parsed[i], _ = strconv.Atoi(part)
	}
	return parsed
}

func compareVersions(a, b []int) int {
	for i := 0; i < len(a) || i < len(b); i++ {
		var av, bv int
		if i < len(a) {
			av = a[i]
		}
		if i < len(b) {
			bv = b[i]
		}
		if av != bv {
			return av - bv
		}
	}
	return 0
}

func groupExtensionVersions(entries []os.DirEntry) map[string][]extensionVersion {
	groups := make(map[string][]extensionVersion)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		match := vscodeExtensionDirPattern.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}

		id, version, platform := match[1], match[2], match[3]
		key := id + platform
		groups[key] = append(groups[key], extensionVersion{
			dirName: entry.Name(),
			version: parseVersion(version),
		})
	}

	return groups
}

func staleExtensionDirs(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var stale []string
	for _, versions := range groupExtensionVersions(entries) {
		if len(versions) < 2 {
			continue
		}

		sort.Slice(versions, func(i, j int) bool {
			return compareVersions(versions[i].version, versions[j].version) > 0
		})

		for _, outdated := range versions[1:] {
			stale = append(stale, outdated.dirName)
		}
	}

	return stale, nil
}

func (v VSCodeExtensions) Size() (int64, bool) {
	root, err := v.path()
	if err != nil {
		return 0, false
	}

	stale, err := staleExtensionDirs(root)
	if err != nil {
		return 0, false
	}

	var total int64
	for _, name := range stale {
		if size, ok := dirSize(filepath.Join(root, name)); ok {
			total += size
		}
	}

	return total, true
}

func (v VSCodeExtensions) Clean() {
	root, err := v.path()
	if err != nil {
		fmt.Println("failed to resolve home directory:", err)
		return
	}

	stale, err := staleExtensionDirs(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s not found (nothing to clean)\n", root)
			return
		}
		fmt.Printf("failed to read %s: %v\n", root, err)
		return
	}

	if len(stale) == 0 {
		fmt.Println("no outdated vscode extensions found")
		return
	}

	for _, name := range stale {
		path := filepath.Join(root, name)
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("failed to remove %s: %v\n", path, err)
			continue
		}
		fmt.Printf("removed %s\n", path)
	}
}

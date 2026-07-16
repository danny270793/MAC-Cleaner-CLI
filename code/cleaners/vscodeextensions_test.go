package cleaners

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseVersion(t *testing.T) {
	got := parseVersion("2.1.211")
	want := []int{2, 1, 211}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseVersion(\"2.1.211\") = %v, want %v", got, want)
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		name string
		a, b []int
		want int
	}{
		{"greater", []int{2, 1, 211}, []int{2, 1, 9}, 1},
		{"lesser", []int{2, 1, 9}, []int{2, 1, 211}, -1},
		{"equal with trailing zero", []int{1, 0}, []int{1, 0, 0}, 0},
		{"equal", []int{1, 2}, []int{1, 2}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := compareVersions(c.a, c.b)
			gotSign := sign(got)
			wantSign := sign(c.want)
			if gotSign != wantSign {
				t.Fatalf("compareVersions(%v, %v) = %d, want sign %d", c.a, c.b, got, wantSign)
			}
		})
	}
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	default:
		return 0
	}
}

func TestStaleExtensionDirsKeepsLatestPerExtension(t *testing.T) {
	root := t.TempDir()
	dirs := []string{
		"anthropic.claude-code-2.1.211-darwin-arm64",
		"anthropic.claude-code-2.1.210-darwin-arm64",
		"anthropic.claude-code-2.1.9-darwin-arm64",
		"ms-python.python-2024.10.1",
		"ms-python.python-2024.2.0",
		"golang.go-0.42.0",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatalf("failed to create %s: %v", d, err)
		}
	}

	stale, err := staleExtensionDirs(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	staleSet := make(map[string]bool)
	for _, s := range stale {
		staleSet[s] = true
	}

	wantStale := []string{
		"anthropic.claude-code-2.1.210-darwin-arm64",
		"anthropic.claude-code-2.1.9-darwin-arm64",
		"ms-python.python-2024.2.0",
	}
	for _, w := range wantStale {
		if !staleSet[w] {
			t.Errorf("expected %s to be marked stale", w)
		}
	}

	wantKept := []string{
		"anthropic.claude-code-2.1.211-darwin-arm64",
		"ms-python.python-2024.10.1",
		"golang.go-0.42.0",
	}
	for _, k := range wantKept {
		if staleSet[k] {
			t.Errorf("expected %s to be kept, but it was marked stale", k)
		}
	}

	if len(stale) != len(wantStale) {
		t.Errorf("expected %d stale dirs, got %d: %v", len(wantStale), len(stale), stale)
	}
}

func TestVSCodeExtensionsCleanRemovesOnlyStale(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	ext := filepath.Join(home, ".vscode", "extensions")
	dirs := []string{
		"anthropic.claude-code-2.1.211-darwin-arm64",
		"anthropic.claude-code-2.1.210-darwin-arm64",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(ext, d), 0o755); err != nil {
			t.Fatalf("failed to create %s: %v", d, err)
		}
	}

	VSCodeExtensions{}.Clean()

	if _, err := os.Stat(filepath.Join(ext, "anthropic.claude-code-2.1.211-darwin-arm64")); err != nil {
		t.Fatalf("expected latest version to remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(ext, "anthropic.claude-code-2.1.210-darwin-arm64")); !os.IsNotExist(err) {
		t.Fatalf("expected older version to be removed, err=%v", err)
	}
}

func TestVSCodeExtensionsNameAndMissingSize(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	v := VSCodeExtensions{}
	if v.Name() != "vscode extensions" {
		t.Fatalf("unexpected name %q", v.Name())
	}

	size, ok := v.Size()
	if ok {
		t.Fatalf("expected ok=false, got size=%d", size)
	}
}

package cleaners

import "testing"

func TestDockerName(t *testing.T) {
	if (Docker{}).Name() != "docker" {
		t.Fatalf("unexpected name %q", (Docker{}).Name())
	}
}

func TestDockerSizeIsUnmeasurable(t *testing.T) {
	size, ok := (Docker{}).Size()
	if ok || size != 0 {
		t.Fatalf("expected size=0, ok=false, got size=%d ok=%v", size, ok)
	}
}

func TestDockerCleanWhenNotInstalled(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	output := captureStdout(t, func() {
		Docker{}.Clean()
	})

	want := "docker not installed (nothing to clean)\n"
	if output != want {
		t.Fatalf("expected output %q, got %q", want, output)
	}
}

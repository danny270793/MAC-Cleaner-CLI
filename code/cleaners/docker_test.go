package cleaners

import "testing"

func TestDockerName(t *testing.T) {
	if (Docker{}).Name() != "docker" {
		t.Fatalf("unexpected name %q", (Docker{}).Name())
	}
}

func TestDockerSizeWhenNotInstalled(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	size, ok := (Docker{}).Size()
	if ok || size != 0 {
		t.Fatalf("expected size=0, ok=false, got size=%d ok=%v", size, ok)
	}
}

func TestDockerCleanWhenNotInstalled(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	var size int64
	var ok bool
	output := captureStdout(t, func() {
		size, ok = Docker{}.Clean()
	})

	want := "docker not installed (nothing to clean)\n"
	if output != want {
		t.Fatalf("expected output %q, got %q", want, output)
	}
	if ok || size != 0 {
		t.Fatalf("expected size=0, ok=false, got size=%d ok=%v", size, ok)
	}
}

func TestParseDockerReclaimedSpace(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   int64
		wantOk bool
	}{
		{
			name: "gigabytes",
			output: `Deleted build cache objects:
bqvq0lzdchqwvnyrhdykr3su4

Total reclaimed space: 2.521GB
`,
			want:   2521000000,
			wantOk: true,
		},
		{
			name:   "megabytes",
			output: "Total reclaimed space: 515.9MB\n",
			want:   515900000,
			wantOk: true,
		},
		{
			name:   "zero bytes",
			output: "Total reclaimed space: 0B\n",
			want:   0,
			wantOk: true,
		},
		{
			name:   "no matching line",
			output: "nothing to clean\n",
			want:   0,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseDockerReclaimedSpace(tt.output)
			if ok != tt.wantOk {
				t.Fatalf("expected ok=%v, got ok=%v", tt.wantOk, ok)
			}
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestParseHumanSize(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   int64
		wantOk bool
	}{
		{name: "gigabytes with percentage suffix", input: "2.1GB (65%)", want: 2100000000, wantOk: true},
		{name: "zero bytes with percentage suffix", input: "0B (0%)", want: 0, wantOk: true},
		{name: "kilobytes", input: "10.75kB", want: 10750, wantOk: true},
		{name: "not a size", input: "N/A", want: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseHumanSize(tt.input)
			if ok != tt.wantOk {
				t.Fatalf("expected ok=%v, got ok=%v", tt.wantOk, ok)
			}
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

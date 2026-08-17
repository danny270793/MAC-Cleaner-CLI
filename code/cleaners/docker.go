package cleaners

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Docker struct{}

func (Docker) Name() string {
	return "docker"
}

func (Docker) Size() (int64, bool) {
	if _, err := exec.LookPath("docker"); err != nil {
		return 0, false
	}

	output, err := exec.Command("docker", "system", "df", "--format", "{{.Reclaimable}}").Output()
	if err != nil {
		return 0, false
	}

	var total int64
	found := false
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if size, ok := parseHumanSize(line); ok {
			total += size
			found = true
		}
	}

	return total, found
}

func (Docker) Clean() (int64, bool) {
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("docker not installed (nothing to clean)")
		return 0, false
	}

	var captured bytes.Buffer
	cmd := exec.Command("docker", "system", "prune", "--all", "--volumes", "--force")
	cmd.Stdout = io.MultiWriter(os.Stdout, &captured)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run docker system prune:", err)
		return 0, false
	}

	return parseDockerReclaimedSpace(captured.String())
}

var (
	dockerReclaimedSpacePattern = regexp.MustCompile(`(?i)Total reclaimed space:\s*(.+)`)
	humanSizePattern            = regexp.MustCompile(`(?i)^\s*([\d.]+)\s*([a-zA-Z]*B)`)
)

func parseDockerReclaimedSpace(output string) (int64, bool) {
	match := dockerReclaimedSpacePattern.FindStringSubmatch(output)
	if match == nil {
		return 0, false
	}

	return parseHumanSize(match[1])
}

func parseHumanSize(s string) (int64, bool) {
	match := humanSizePattern.FindStringSubmatch(s)
	if match == nil {
		return 0, false
	}

	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, false
	}

	units := []string{"b", "kb", "mb", "gb", "tb", "pb", "eb", "zb", "yb"}
	unit := strings.ToLower(match[2])

	multiplier := 1.0
	for _, u := range units {
		if unit == u {
			return int64(value * multiplier), true
		}
		multiplier *= 1000
	}

	return 0, false
}

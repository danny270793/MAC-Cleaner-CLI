package cleaners

import (
	"fmt"
	"os"
	"os/exec"
)

type Docker struct{}

func (Docker) Name() string {
	return "docker"
}

func (Docker) Clean() {
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("docker not installed (nothing to clean)")
		return
	}

	cmd := exec.Command("docker", "system", "prune")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run docker system prune:", err)
	}
}

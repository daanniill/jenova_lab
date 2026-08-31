package lab

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		fmt.Print("docker up")

	case "status":
		fmt.Print("docker status")

	case "down":
		fmt.Print("docker down")

	default:
		os.Exit(1)
	}
}

func runDocker(args ...string) {
	cmd := exec.Command("docker", args...)

	// mapping streams directly to global stream
	// allows us to interact with program
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("docker command failed: %v\n", err)
		os.Exit(1)
	}
}
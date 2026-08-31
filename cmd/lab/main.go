package lab

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const composeFile = "compose.generated.yaml"

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

func upCommand(args []string) {
	flags := flag.NewFlagSet("up", flag.ExitOnError)

	services := flags.Int("services", 3, "number of services")

	flags.Parse(args)

	if *services < 1 {
		fmt.Println("services must be at least 1")
		os.Exit(1)
	}

	fmt.Printf("generate a compose file")

	fmt.Printf("Starting %d-service chain...\n", *services)

	runDocker("compose", "-f", composeFile, "up", "--build", "-d")
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
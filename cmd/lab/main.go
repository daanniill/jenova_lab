package lab

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const composeFile = "compose.generated.yaml"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		upCommand(os.Args[2:])
	
	//lists containers that are running
	case "status":
		runDocker("compose", "-f", composeFile, "ps")

	case "down":
		runDocker("compose", "-f", composeFile, "down")

	default:
		usage()
		os.Exit(1)
	}
}

func upCommand(args []string) {
	// create a new flag set and define services flag that tracks how many services should be launched
	flags := flag.NewFlagSet("up", flag.ExitOnError)
	services := flags.Int("services", 3, "number of services")
	flags.Parse(args)

	if *services < 1 {
		fmt.Println("services must be at least 1")
		os.Exit(1)
	}

	// generate a compose file
	if err := generateCompose(*services); err != nil {
		fmt.Printf("failed to generate compose file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting %d-service chain...\n", *services)

	runDocker("compose", "-f", composeFile, "up", "--build", "-d")
}

func generateCompose(count int) error {
	var b strings.Builder

	b.WriteString("services:\n")

	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("service-%d", i)
		port := 8000 + i
		
		fmt.Fprintf(&b, "	%s:\n", name)

		b.WriteString("		build:\n")
		b.WriteString("			context: .\n")
		b.WriteString("     dockerfile: service/Dockerfile\n")

		b.WriteString("		environment:\n")
		fmt.Fprintf(&b, "			SERVICE_NAME: %s\n", name)
		b.WriteString("			PORT: 8080\n")

		// Every service except the last points to the next service.
		if i < count {
			next := fmt.Sprintf("service-%d", i+1)
			fmt.Fprintf(&b,"			DOWNSTREAM_URL: http://%s:8080\n",next)
		}

		b.WriteString("    ports:\n")
		fmt.Fprintf(&b, "      - \"%d:8080\"\n", port)
	}

	// 0644 grants read/write permissions to owner, and read-only to others
	return os.WriteFile(composeFile, []byte(b.String()), 0644)
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

func usage() {
	fmt.Println(`Usage:

  lab up [--services N]
  lab status
  lab down`)
}
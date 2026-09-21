package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	
	case "kill":
		if len(os.Args) != 3 {
			fmt.Println("usage: lab kill <service>")
			os.Exit(1)
		}
		killService(os.Args[2])
	
	case "reset":
		resetLab()
	
	case "latency":
    if len(os.Args) != 4 {
        fmt.Println("usage: lab latency <service> <duration>")
        os.Exit(1)
    }
    setLatency(os.Args[2], os.Args[3])

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

func killService(name string) { 
	if _, err := servicePort(name); err != nil {
		fmt.Printf("failed to process service name: %v\n", err)
		os.Exit(1)
	}

	runDocker("compose", "-f", "compose.generated.yaml", "kill", name)
}

func setLatency(name string, rawDuration string) {
	port, err := servicePort(name)
	if err != nil {
		fmt.Printf("failed to process service name: %v\n", err)
		os.Exit(1)
	}

	duration, err := time.ParseDuration(rawDuration)
	if err != nil {
		fmt.Printf("invalid duration %q, %v; try 500ms, 1s, 1500ms", rawDuration, err)
		os.Exit(1)
	}

	if duration < 0 {
		fmt.Printf("duration must be positive")
		os.Exit(1)
	}

	// Send post request to service to apply latency
	endpoint := fmt.Sprintf("http://localhost:%d/fault/latency?duration=%s", port, url.QueryEscape(duration.String()))

	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		fmt.Printf("post request failed: %v", err)
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("could not reach %s: %v", name, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s returned %s: %s", name, resp.Status, strings.TrimSpace(string(body)))
	}

	fmt.Printf("%s: latency = %s\n", name, duration)
}

// Recreate everything.
	//
	// This:
	//   1. starts killed containers again
	//   2. clears injected in-memory latency
func resetLab() {
	runDocker("compose", "-f", "compose.generated.yaml", "up", "-d", "--force-recreate")
}

func generateCompose(count int) error {
	var b strings.Builder

	b.WriteString("services:\n")

	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("service-%d", i)
		port := 8000 + i

		fmt.Fprintf(&b, "  %s:\n", name)

		b.WriteString("    build:\n")
		b.WriteString("       context: .\n")
		b.WriteString("       dockerfile: service/Dockerfile\n\n")

		b.WriteString("    environment:\n")
		fmt.Fprintf(&b, "       SERVICE_NAME: %s\n", name)
		b.WriteString("       PORT: 8080\n")

		// Every service except the last points to the next service.
		if i < count {
			next := fmt.Sprintf("service-%d", i+1)
			fmt.Fprintf(&b, "       DOWNSTREAM_URL: http://%s:8080\n\n", next)
		}

		b.WriteString("    ports:\n")
		fmt.Fprintf(&b, "       - \"%d:8080\"\n", port)
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

// service-x -> localhost:800x
func servicePort(name string) (int, error) {
	const prefix = "service-"

	if !strings.HasPrefix(name, prefix) {
		return 0, fmt.Errorf("invalid service name %q", name)
	}

	numberString := strings.TrimPrefix(name, prefix)

	number, err := strconv.Atoi(numberString)
	if err != nil || number < 1 {
		return 0, fmt.Errorf("invalid service name %q", name)
	}

	return 8000 + number, nil
}

func usage() {
	fmt.Println(`Usage:

  lab up [--services N]
  lab status
  lab down`)
}

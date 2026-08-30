package lab

import (
	"fmt"
	"os"
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
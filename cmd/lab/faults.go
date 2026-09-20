package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
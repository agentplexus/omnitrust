package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/posture/server"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

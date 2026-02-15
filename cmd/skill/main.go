package main

import (
	"fmt"
	"os"

	"github.com/silvinalucero/skill_cli/internal/cli"
)

// Build information. Populated at build-time by GoReleaser.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version information
	cli.SetVersion(version, commit, date)

	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/edgardnogueira/swagger-to-http-file/internal/infrastructure/cli"
)

// Version information - will be set during build with ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info
	cli.SetVersionInfo(version, commit, date)
	
	// Execute CLI
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

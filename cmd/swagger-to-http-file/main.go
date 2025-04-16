package main

import (
	"fmt"
	"os"

	"github.com/edgardnogueira/swagger-to-http-file/internal/infrastructure/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

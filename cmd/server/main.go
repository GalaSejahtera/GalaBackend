package main

import (
	"fmt"
	"os"

	cmd "safeworkout/pkg/cmd/server"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

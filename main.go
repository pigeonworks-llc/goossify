// Package main is the entry point for the goossify CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/pigeonworks-llc/goossify/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/melkeydev/go-blueprint/cmd"
)

func main() {
	if err := cmd.New().Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

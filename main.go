package main

import (
	"fmt"
	"os"

	"go.aponeill.com/fill/pkg/cmd/fill"
)

func main() {
	cmd := fill.NewCommand()
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

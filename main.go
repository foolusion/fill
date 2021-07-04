package main

import (
	"fmt"
	"os"

	"aponeill.com/fill/pkg/cmd/fill"
)

func main() {
	cmd := fill.NewCommand()
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

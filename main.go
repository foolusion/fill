package main

import (
	"embed"
	"fmt"
	"os"

	"go.aponeill.com/fill/pkg/cmd/fill"
)

//go:embed res/*
var fs embed.FS

func main() {
	cmd := fill.NewCommand(fs)
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

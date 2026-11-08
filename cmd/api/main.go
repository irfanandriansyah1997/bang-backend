package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(w io.Writer) error {
	_, err := fmt.Fprintln(w, startupMessage("Benah API"))

	return err
}

func startupMessage(appName string) string {
	if appName != "" {
		return appName + " starting"
	}

	return ""
}

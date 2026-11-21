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
	if _, err := fmt.Fprintln(w, startupMessage("Benah API")); err != nil {
		return fmt.Errorf("write startup message: %w", err)
	}

	return nil
}

func startupMessage(appName string) string {
	if appName != "" {
		return appName + " starting"
	}

	return ""
}

package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
)

func main() {
	initApp()
	initAppHelp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

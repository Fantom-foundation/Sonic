package main

import (
	"fmt"
	"os"

	"github.com/Fantom-foundation/go-opera/cmd/sonictool/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

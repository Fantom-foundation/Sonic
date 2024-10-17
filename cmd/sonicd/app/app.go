package app

import (
	"os"
)

func Run() error {
	initApp()
	initAppHelp()
	return app.Run(os.Args)
}

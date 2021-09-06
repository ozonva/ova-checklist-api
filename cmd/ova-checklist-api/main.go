package main

import (
	"os"

	"github.com/ozonva/ova-checklist-api/internal/application"
)

func main() {
	defer func() {
		const exitCodeOnFail = 1
		if err := recover(); err != nil {
			os.Exit(exitCodeOnFail)
		}
	}()
	application.Run()
}

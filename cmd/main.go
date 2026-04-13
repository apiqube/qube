package main

import (
	"os"

	"github.com/apiqube/qube/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

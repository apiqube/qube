// Package main is the entry point for the qube CLI binary.
package main

import (
	"context"

	"github.com/charmbracelet/fang"

	"github.com/apiqube/qube/internal/commands"
)

func main() {
	// fang wraps cobra and styles help screens, errors, and the version block.
	// It writes its own styled output on error, so we don't need a custom
	// fmt.Fprintln fallback here.
	_ = fang.Execute(context.Background(), commands.Root())
}

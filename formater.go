package main

import (
	"io"

	"github.com/arl/gitstatus"
)

// A formater writes the status of a Git working tree in a given format.
type formater interface {
	// Format writes the representation of a git status.
	Format(io.Writer, *gitstatus.Status) error
}

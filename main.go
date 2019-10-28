package main

import (
	"errors"
	"fmt"

	"github.com/arl/gitstatus"
	"github.com/arl/gitstatus/format/json"
	"github.com/arl/gitstatus/format/tmux"
)

var errUnknownOutputFormat = errors.New("unknown output format")

func main() {
	// parse cli options.
	dir, format, quiet, cfg := parseOptions()

	// handle directory change.
	if dir != "." {
		popDir, err := pushdir(dir)
		check(err, quiet)
		defer func() {
			check(popDir(), quiet)
		}()
	}

	// retrieve git status.
	st, err := gitstatus.New()
	check(err, quiet)

	// register formaters
	formaters := make(map[string]gitstatus.Formater)
	formaters["json"] = &json.Formater{}
	formaters["tmux"] = &tmux.Formater{Config: cfg.Tmux}

	formater, ok := formaters[format]
	if !ok {
		check(errUnknownOutputFormat, quiet)
	}

	// format and print
	out, err := formater.Format(st)
	check(err, quiet)
	fmt.Print(out)
}

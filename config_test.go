//go:build !windows
// +build !windows

package main

import (
	"flag"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestScripts(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("GetWd error: %v", err)
	}
	params := testscript.Params{
		Dir: "testdata",
		Setup: func(env *testscript.Env) error {
			env.Setenv("GITMUX_DIR", wd)
			return nil
		},
		UpdateScripts: *updateGolden,
		TestWork:      true,
	}
	if err := gotooltest.Setup(&params); err != nil {
		t.Errorf("gotooltest.Setup error: %v", err)
	}
	testscript.Run(t, params)
}

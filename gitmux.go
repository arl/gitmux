package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arl/gitstatus"
	"gopkg.in/yaml.v2"

	"github.com/arl/gitmux/format/json"
	"github.com/arl/gitmux/format/tmux"
)

var version = "<<development version>>"

var usage = `gitmux ` + version + `
Usage: gitmux [options] [dir]

gitmux prints the status of a Git working tree as a tmux format string.
If directory is not given, it default to the working directory.  

Options:
  -cfg cfgfile    use cfgfile when printing git status.
  -printcfg       prints default configuration file.
  -dbg            outputs Git status as JSON and print errors.
  -V              prints gitmux version and exits.
`

// Config configures output formatting.
type Config struct{ Tmux tmux.Config }

var defaultCfg = Config{Tmux: tmux.DefaultCfg}

func parseOptions() (dir string, dbg bool, cfg Config) {
	dbgOpt := flag.Bool("dbg", false, "")
	cfgOpt := flag.String("cfg", "", "")
	printCfgOpt := flag.Bool("printcfg", false, "")
	versionOpt := flag.Bool("V", false, "")
	flag.Bool("q", true, "")   // unused, kept for retrocompatibility.
	flag.String("fmt", "", "") // unused, kept for retrocompatibility.
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()

	dir = "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	if *versionOpt {
		fmt.Println(version)
		os.Exit(0)
	}

	if *printCfgOpt {
		enc := yaml.NewEncoder(os.Stdout)
		check(enc.Encode(&defaultCfg), *dbgOpt)
		enc.Close()
		os.Exit(0)
	}

	cfg = defaultCfg

	if *cfgOpt != "" {
		f, err := os.Open(*cfgOpt)
		check(err, *dbgOpt)

		dec := yaml.NewDecoder(f)
		check(dec.Decode(&cfg), *dbgOpt)
	}

	return dir, *dbgOpt, cfg
}

func pushdir(dir string) (popdir func() error, err error) {
	pwd := ""
	if pwd, err = os.Getwd(); err != nil {
		return nil, err
	}

	if err = os.Chdir(dir); err != nil {
		return nil, err
	}

	return func() error { return os.Chdir(pwd) }, nil
}

func check(err error, dbg bool) {
	if err == nil {
		return
	}
	if dbg {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	os.Exit(1)
}

func main() {
	dir, dbg, cfg := parseOptions()

	// handle directory change.
	if dir != "." {
		popDir, err := pushdir(dir)
		check(err, dbg)
		defer func() {
			check(popDir(), dbg)
		}()
	}

	// retrieve git status.
	st, err := gitstatus.New()
	check(err, dbg)

	// select defauit formater
	var formater formater = &tmux.Formater{Config: cfg.Tmux}
	if dbg {
		formater = &json.Formater{}
	}

	// format and print
	check(formater.Format(os.Stdout, st), dbg)
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/arl/gitmux/format/json"
	"github.com/arl/gitmux/format/tmux"
	"github.com/arl/gitstatus"
	"gopkg.in/yaml.v2"
)

func check(err error, dbg bool) {
	if err != nil && dbg {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// Config configures output formatting.
type Config struct{ Tmux tmux.Config }

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

var defaultCfg = Config{Tmux: tmux.DefaultCfg}

func parseOptions() (dir string, dbg bool, cfg Config) {
	dbgOpt := flag.Bool("dbg", false, "")
	cfgOpt := flag.String("cfg", "", "")
	printCfgOpt := flag.Bool("printcfg", false, "")
	versionOpt := flag.Bool("V", false, "")
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()
	dir = "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	cfg = defaultCfg
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
	if *cfgOpt != "" {
		f, err := os.Open(*cfgOpt)
		check(err, *dbgOpt)
		dec := yaml.NewDecoder(f)
		check(dec.Decode(&cfg), *dbgOpt)
	}
	return dir, *dbgOpt, cfg
}

type popdir func() error

func pushdir(dir string) (popdir, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(dir)
	if err != nil {
		return nil, err
	}

	return func() error { return os.Chdir(pwd) }, nil
}

var errUnknownOutputFormat = errors.New("unknown output format")

func main() {
	// parse cli options.
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

	// select formater
	var formater formater = &tmux.Formater{Config: cfg.Tmux}
	if dbg {
		formater = &json.Formater{}
	}

	// format and print
	err = formater.Format(os.Stdout, st)
	check(err, dbg)
}

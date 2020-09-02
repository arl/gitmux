package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/arl/gitstatus"
	"gopkg.in/yaml.v2"

	"github.com/arl/gitmux/format/json"
	"github.com/arl/gitmux/format/tmux"
)

var version = "<<development version>>"

var _usage = `gitmux ` + version + `
Usage: gitmux [options] [dir]

gitmux prints the status of a Git working tree as a tmux format string.
If directory is not given, it default to the working directory.  

Options:
  -cfg cfgfile    use cfgfile when printing git status.
  -printcfg       prints default configuration file.
  -dbg            outputs Git status as JSON and print errors.
  -timeout DUR    exits if still running after given duration (ex: 2s, 500ms).
  -V              prints gitmux version and exits.
`

// Config configures output formatting.
type Config struct{ Tmux tmux.Config }

var _defaultCfg = Config{Tmux: tmux.DefaultCfg}

// duration is time.Duration usable as command line flag.
type duration time.Duration

func (d duration) String() string {
	if d == 0 {
		return "none"
	}

	return time.Duration(d).String()
}

func (d *duration) Set(s string) error {
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = duration(dur)
	return nil
}

func parseOptions() (ctx context.Context, dir string, dbg bool, cfg Config) {
	dbgOpt := flag.Bool("dbg", false, "")
	cfgOpt := flag.String("cfg", "", "")
	printCfgOpt := flag.Bool("printcfg", false, "")
	versionOpt := flag.Bool("V", false, "")
	timeout := duration(0)
	flag.Var(&timeout, "timeout", "")

	flag.Bool("q", true, "")   // unused, kept for retro-compatibility.
	flag.String("fmt", "", "") // unused, kept for retro-compatibility.
	flag.Usage = func() {
		fmt.Println(_usage)
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
		check(enc.Encode(&_defaultCfg), *dbgOpt)
		enc.Close()
		os.Exit(0)
	}

	cfg = _defaultCfg

	if *cfgOpt != "" {
		f, err := os.Open(*cfgOpt)
		check(err, *dbgOpt)

		dec := yaml.NewDecoder(f)
		check(dec.Decode(&cfg), *dbgOpt)
	}

	ctx = context.Background()
	if timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, time.Duration(timeout))
	}

	return ctx, dir, *dbgOpt, cfg
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
	ctx, dir, dbg, cfg := parseOptions()

	// handle directory change.
	if dir != "." {
		popDir, err := pushdir(dir)

		check(err, dbg)
		defer func() {
			check(popDir(), dbg)
		}()
	}

	// retrieve git status.
	st, err := gitstatus.NewWithContext(ctx)
	check(err, dbg)

	// select defauit formater
	var formater formater = &tmux.Formater{Config: cfg.Tmux}
	if dbg {
		formater = &json.Formater{}
	}

	// format and print
	check(formater.Format(os.Stdout, st), dbg)
}

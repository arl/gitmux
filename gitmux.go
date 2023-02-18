package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/arl/gitstatus"
	"gopkg.in/yaml.v3"

	"github.com/arl/gitmux/json"
	"github.com/arl/gitmux/tmux"
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

func parseOptions() (ctx context.Context, cancel func(), dir string, dbg bool, cfg Config) {
	var (
		dbgOpt      = flag.Bool("dbg", false, "")
		cfgOpt      = flag.String("cfg", "", "")
		printCfgOpt = flag.Bool("printcfg", false, "")
		versionOpt  = flag.Bool("V", false, "")
		timeoutOpt  = flag.Duration("timeout", 0, "")
	)

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
		os.Stdout.Write(cfgBytes)
		os.Exit(0)
	}

	cfg = defaultCfg

	if *cfgOpt != "" {
		f, err := os.Open(*cfgOpt)
		check(err, *dbgOpt)

		dec := yaml.NewDecoder(f)
		check(dec.Decode(&cfg), *dbgOpt)
	}

	if *timeoutOpt != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), *timeoutOpt)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	return ctx, cancel, dir, *dbgOpt, cfg
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
	ctx, cancel, dir, dbg, cfg := parseOptions()
	defer cancel()

	// Handle directory change.
	if dir != "." {
		popDir, err := pushdir(dir)

		check(err, dbg)
		defer func() {
			check(popDir(), dbg)
		}()
	}

	// Retrieve git status.
	st, err := gitstatus.NewWithContext(ctx)
	check(err, dbg)

	// Interface that writes a particular representation of a gitstatus.Status
	type formater interface {
		Format(io.Writer, *gitstatus.Status) error
	}

	// Set defauit formater.
	var fmter formater = &tmux.Formater{Config: cfg.Tmux}
	if dbg {
		fmter = &json.Formater{}
	}

	check(fmter.Format(os.Stdout, st), dbg)
}

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

func check(err error, quiet bool) {
	if err != nil {
		if !quiet {
			fmt.Println("error:", err)
		}
		os.Exit(1)
	}
}

// Config configures output formatting.
type Config struct{ Tmux tmux.Config }

var version = "<<development version>>"

var usage = `gitmux ` + version + `
Usage: gitmux [options] [dir]

gitmux prints the status of a Git working tree.
If directory is not given, it default to the working directory.  

Options:
  -q              be quiet. In case of errors, don't print nothing.
  -fmt            output format, defaults to json.
      json        prints status as a JSON object.
      tmux        prints status as a tmux format string.
  -cfg cfgfile    use cfgfile when printing git status.
  -printcfg       prints default configuration file.
  -V		  prints gitmux version and exits.
`

var defaultCfg = Config{Tmux: tmux.DefaultCfg}

func parseOptions() (dir string, format string, quiet bool, cfg Config) {
	fmtOpt := flag.String("fmt", "json", "")
	cfgOpt := flag.String("cfg", "", "")
	printCfgOpt := flag.Bool("printcfg", false, "")
	quietOpt := flag.Bool("q", false, "")
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
		check(enc.Encode(&defaultCfg), *quietOpt)
		enc.Close()
		os.Exit(0)
	}
	if *cfgOpt != "" {
		f, err := os.Open(*cfgOpt)
		check(err, *quietOpt)
		dec := yaml.NewDecoder(f)
		check(dec.Decode(&cfg), *quietOpt)
	}
	return dir, *fmtOpt, *quietOpt, cfg
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
	formaters := make(map[string]formater)
	formaters["json"] = &json.Formater{}
	formaters["tmux"] = &tmux.Formater{Config: cfg.Tmux}

	formater, ok := formaters[format]
	if !ok {
		check(errUnknownOutputFormat, quiet)
	}

	// format and print
	err = formater.Format(os.Stdout, st)
	check(err, quiet)
}

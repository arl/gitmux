package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arl/gitstatus/format/tmux"
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

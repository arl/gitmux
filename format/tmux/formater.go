package tmux

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/arl/gitstatus"
)

const clear string = "#[fg=default]"

// Config is the configuration of the Git status tmux formatter.
type Config struct {
	// Symbols contains the symbols printed before the Git status components.
	Symbols symbols
	// Styles contains the tmux style strings for symbols and Git status
	// components.
	Styles  styles
	Display display
}

type symbols struct {
	Branch     string // Branch is the string shown before local branch name.
	HashPrefix string // HasPrefix is the string shown before a SHA1 ref.

	Ahead  string // Ahead is the string shown before the ahead count for the local/upstream branch divergence.
	Behind string // Behind is the string shown before the behind count for the local/upstream branch divergence.

	Staged    string // Staged is the string shown before the count of staged files.
	Conflict  string // Conflict is the string shown before the count of files with conflicts.
	Modified  string // Modified is the string shown before the count of modified files.
	Untracked string // Untracked is the string shown before the count of untracked files.
	Stashed   string // Stashed is the string shown before the count of stash entries.
	Clean     string // Clean is the string shown when the working tree is clean.
	Delimiter string // Delimiter is the string shown before flags
}

type styles struct {
	State     string // State is the style string printed before eventual special state.
	Branch    string // Branch is the style string printed before the local branch.
	Remote    string // Remote is the style string printed before the upstream branch.
	Staged    string // Staged is the style string printed before the staged files count.
	Conflict  string // Conflict is the style string printed before the conflict count.
	Modified  string // Modified is the style string printed before the modified files count.
	Untracked string // Untracked is the style string printed before the untracked files count.
	Stashed   string // Stashed is the style string printed before the stash entries count.
	Clean     string // Clean is the style string printed before the clean symbols.
}

type display struct {
	Output []string // Display order and visiblity of options
}

var DefaultCfg = Config{
	Symbols: symbols{
		Branch:     "⎇ ",
		Staged:     "●",
		Conflict:   "✖ ",
		Modified:   "✚ ",
		Untracked:  "…",
		Stashed:    "⚑ ",
		Clean:      "✔",
		Ahead:      "↑·",
		Behind:     "↓·",
		HashPrefix: ":",
		Delimiter:  " - ",
	},
	Styles: styles{
		State:     "#[fg=red,bold]",
		Branch:    "#[fg=white,bold]",
		Remote:    "#[fg=cyan]",
		Staged:    "#[fg=green,bold]",
		Conflict:  "#[fg=red,bold]",
		Modified:  "#[fg=red,bold]",
		Untracked: "#[fg=magenta,bold]",
		Stashed:   "#[fg=cyan,bold]",
		Clean:     "#[fg=green,bold]",
	},
	Display: display{
		Output: []string{"branch", "remote", "flags"},
	},
}

// A Formater formats git status to a tmux style string.
type Formater struct {
	Config
	b  bytes.Buffer
	st *gitstatus.Status
}

// Format writes st as json into w.
func (f *Formater) Format(w io.Writer, st *gitstatus.Status) error {
	f.st = st
	f.clear()

	// overall working tree state
	if f.st.IsInitial {
		fmt.Fprintf(w, "%s%s [no commits yet]", f.Styles.Branch, f.st.LocalBranch)
		f.flags()
	} else {
		for _, order := range f.Display.Output {
			switch order {
			case "branch":
				f.specialState()
			case "remote":
				f.remote()
			case "flags":
				f.flags()
			}
		}
	}

	_, err := f.b.WriteTo(w)

	return err
}

func (f *Formater) specialState() {
	f.clear()
	switch f.st.State {
	case gitstatus.Rebasing:
		fmt.Fprintf(&f.b, "%s[rebase] ", f.Styles.State)
	case gitstatus.AM:
		fmt.Fprintf(&f.b, "%s[am] ", f.Styles.State)
	case gitstatus.AMRebase:
		fmt.Fprintf(&f.b, "%s[am-rebase] ", f.Styles.State)
	case gitstatus.Merging:
		fmt.Fprintf(&f.b, "%s[merge] ", f.Styles.State)
	case gitstatus.CherryPicking:
		fmt.Fprintf(&f.b, "%s[cherry-pick] ", f.Styles.State)
	case gitstatus.Reverting:
		fmt.Fprintf(&f.b, "%s[revert] ", f.Styles.State)
	case gitstatus.Bisecting:
		fmt.Fprintf(&f.b, "%s[bisect] ", f.Styles.State)
	case gitstatus.Default:
		fmt.Fprintf(&f.b, "%s%s", f.Styles.Branch, f.Symbols.Branch)
	}
	f.currentRef()
}

func (f *Formater) remote() {
	f.clear()
	if f.st.RemoteBranch != "" {
		fmt.Fprintf(&f.b, "..%s%s", f.Styles.Remote, f.st.RemoteBranch)
		f.divergence()
	}
}

func (f *Formater) clear() {
	// clear global style
	f.b.WriteString(clear)
}

func (f *Formater) currentRef() {
	f.clear()
	if f.st.IsDetached {
		fmt.Fprintf(&f.b, "%s%s", f.Symbols.HashPrefix, f.st.HEAD)
		return
	}

	fmt.Fprintf(&f.b, "%s", f.st.LocalBranch)
}

func (f *Formater) divergence() {
	f.clear()

	pref := " "
	if f.st.BehindCount != 0 {
		fmt.Fprintf(&f.b, " %s%d", f.Symbols.Behind, f.st.BehindCount)
		pref = ""
	}

	if f.st.AheadCount != 0 {
		fmt.Fprintf(&f.b, "%s%s%d", pref, f.Symbols.Ahead, f.st.AheadCount)
	}
}

func (f *Formater) flags() {
	f.clear()
	f.b.WriteString(f.Symbols.Delimiter)

	if f.st.IsClean {
		fmt.Fprintf(&f.b, "%s%s", f.Styles.Clean, f.Symbols.Clean)
		return
	}

	var flags []string

	if f.st.NumStaged != 0 {
		flags = append(flags,
			fmt.Sprintf("%s%s%d", f.Styles.Staged, f.Symbols.Staged, f.st.NumStaged))
	}

	if f.st.NumConflicts != 0 {
		flags = append(flags,
			fmt.Sprintf("%s%s%d", f.Styles.Conflict, f.Symbols.Conflict, f.st.NumConflicts))
	}

	if f.st.NumModified != 0 {
		flags = append(flags,
			fmt.Sprintf("%s%s%d", f.Styles.Modified, f.Symbols.Modified, f.st.NumModified))
	}

	if f.st.NumStashed != 0 {
		flags = append(flags,
			fmt.Sprintf("%s%s%d", f.Styles.Stashed, f.Symbols.Stashed, f.st.NumStashed))
	}

	if f.st.NumUntracked != 0 {
		flags = append(flags,
			fmt.Sprintf("%s%s%d", f.Styles.Untracked, f.Symbols.Untracked, f.st.NumUntracked))
	}

	f.b.WriteString(strings.Join(flags, " "))
}

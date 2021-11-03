package tmux

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/arl/gitstatus"
	"gopkg.in/yaml.v3"
)

const truncateSymbol string = "..."

// Config is the configuration of the Git status tmux formatter.
type Config struct {
	// Symbols contains the symbols printed before the Git status components.
	Symbols symbols
	// Styles contains the tmux style strings for symbols and Git status
	// components.
	Styles styles
	// Layout sets the output format of the Git status.
	Layout []string `yaml:",flow"`
	// Options contains additional configuration options.
	Options options
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
}

type styles struct {
	Clear      string // Clear is the style string that clears all styles.
	State      string // State is the style string printed before eventual special state.
	Branch     string // Branch is the style string printed before the local branch.
	Remote     string // Remote is the style string printed before the upstream branch.
	Staged     string // Staged is the style string printed before the staged files count.
	Conflict   string // Conflict is the style string printed before the conflict count.
	Modified   string // Modified is the style string printed before the modified files count.
	Untracked  string // Untracked is the style string printed before the untracked files count.
	Stashed    string // Stashed is the style string printed before the stash entries count.
	Clean      string // Clean is the style string printed before the clean symbols.
	Divergence string // Divergence is the style string printed before divergence count/symbols.
}

const (
	dirLeft  direction = "left"
	dirRight direction = "right"
)

type direction string

func (d *direction) UnmarshalYAML(value *yaml.Node) error {
	s := ""
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("error decoding 'direction': %v", s)
	}
	switch direction(s) {
	case dirLeft:
		*d = dirLeft
	case dirRight:
		*d = dirRight
	default:
		return fmt.Errorf("'direction': unexpected value %v", s)
	}

	// if i, ok := o.value.(int); ok {
	// 	if result, ok := unmarshalerResult[i]; ok {
	// 		return result
	// 	}
	// }
	return nil
}

type options struct {
	// BranchMaxLen is the maximum displayed length for local and remote branch names.
	BranchMaxLen        int       `yaml:"branch_max_len"`
	BranchTrimDirection direction `yaml:"branch_trim_direction"`
}

// DefaultCfg is the default tmux configuration.
var DefaultCfg = Config{
	Symbols: symbols{
		Branch:     "⎇ ",
		Staged:     "● ",
		Conflict:   "✖ ",
		Modified:   "✚ ",
		Untracked:  "… ",
		Stashed:    "⚑ ",
		Clean:      "✔",
		Ahead:      "↑·",
		Behind:     "↓·",
		HashPrefix: ":",
	},
	Styles: styles{
		Clear:      "#[fg=default]",
		State:      "#[fg=red,bold]",
		Branch:     "#[fg=white,bold]",
		Remote:     "#[fg=cyan]",
		Divergence: "#[fg=default]",
		Staged:     "#[fg=green,bold]",
		Conflict:   "#[fg=red,bold]",
		Modified:   "#[fg=red,bold]",
		Untracked:  "#[fg=magenta,bold]",
		Stashed:    "#[fg=cyan,bold]",
		Clean:      "#[fg=green,bold]",
	},
	Layout: []string{"branch", "..", "remote-branch", "divergence", " - ", "flags"},
	Options: options{
		BranchMaxLen:        0,
		BranchTrimDirection: dirRight,
	},
}

// A Formater formats git status to a tmux style string.
type Formater struct {
	Config
	b  bytes.Buffer
	st *gitstatus.Status
}

// Truncates branch name if longer than maxlen. If isremote, the leading
// "<remote>/" is ignored when counting length.
func truncateBranchName(name string, maxlen int, isremote bool, trimDir direction) string {
	remoteName := ""
	branchName := name

	const (
		idxRemote = 0
		idxBranch = 1
		numItems  = 2
	)

	if isremote {
		a := strings.SplitAfterN(name, "/", numItems)
		if len(a) == numItems {
			remoteName = a[idxRemote]
			branchName = a[idxBranch]
		}
	}

	// To count length of characters and extract substring from UTF-8 strings.
	branchNameRune := []rune(branchName)
	truncateSymbolRune := []rune(truncateSymbol)

	if maxlen > 0 && maxlen < len(branchNameRune) {
		nameLen := maxlen - len(truncateSymbolRune)
		if trimDir == dirRight {
			if nameLen > 0 {
				branchName = string(branchNameRune[:nameLen]) + truncateSymbol
			} else {
				branchName = string(truncateSymbolRune[:maxlen])
			}
		} else {
			panic("trimDir: " + trimDir)
		}
	}

	return remoteName + branchName
}

// Format writes st as json into w.
func (f *Formater) Format(w io.Writer, st *gitstatus.Status) error {
	f.st = st
	f.clear()

	// overall working tree state
	if f.st.IsInitial {
		branch := truncateBranchName(f.st.LocalBranch, f.Options.BranchMaxLen, false, f.Options.BranchTrimDirection)
		fmt.Fprintf(w, "%s%s [no commits yet]", f.Styles.Branch, branch)
		f.flags()
		_, err := f.b.WriteTo(w)

		return err
	}

	f.format()
	_, err := f.b.WriteTo(w)

	return err
}

func (f *Formater) format() {
	for _, item := range f.Layout {
		switch item {
		case "branch":
			f.specialState()
		case "remote":
			f.remote()
		case "remote-branch":
			f.remoteBranch()
		case "divergence":
			f.divergence()
		case "flags":
			f.flags()
		default:
			f.clear()
			f.b.WriteString(item)
		}
	}
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
	if f.st.RemoteBranch == "" {
		return
	}

	f.clear()

	branch := truncateBranchName(f.st.RemoteBranch, f.Options.BranchMaxLen, true, f.Options.BranchTrimDirection)
	fmt.Fprintf(&f.b, "%s%s", f.Styles.Remote, branch)

	f.divergence()
}

func (f *Formater) remoteBranch() {
	if f.st.RemoteBranch != "" {
		f.clear()

		branch := truncateBranchName(f.st.RemoteBranch, f.Options.BranchMaxLen, true, f.Options.BranchTrimDirection)
		fmt.Fprintf(&f.b, "%s%s", f.Styles.Remote, branch)

	}
}

func (f *Formater) divergence() {
	if f.st.BehindCount == 0 && f.st.AheadCount == 0 {
		return
	}

	f.clear()
	f.b.WriteByte(' ')
	fmt.Fprintf(&f.b, "%s", f.Styles.Divergence)

	if f.st.BehindCount != 0 {
		fmt.Fprintf(&f.b, "%s%d", f.Symbols.Behind, f.st.BehindCount)
	}

	if f.st.AheadCount != 0 {
		fmt.Fprintf(&f.b, "%s%d", f.Symbols.Ahead, f.st.AheadCount)
	}
}

func (f *Formater) clear() {
	// clear global style
	f.b.WriteString(f.Styles.Clear)
}

func (f *Formater) currentRef() {
	f.clear()

	if f.st.IsDetached {
		fmt.Fprintf(&f.b, "%s%s%s", f.Styles.Branch, f.Symbols.HashPrefix, f.st.HEAD)

		return
	}

	branch := truncateBranchName(f.st.LocalBranch, f.Options.BranchMaxLen, false, f.Options.BranchTrimDirection)
	fmt.Fprintf(&f.b, "%s%s", f.Styles.Branch, branch)
}

func (f *Formater) flags() {
	if f.st.IsClean {
		f.clear()
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

	if len(flags) > 0 {
		f.clear()
		f.b.WriteString(strings.Join(flags, " "))
	}
}

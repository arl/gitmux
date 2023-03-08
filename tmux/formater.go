package tmux

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/arl/gitstatus"
	"gopkg.in/yaml.v3"
)

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

	Insertions string // Insertions is the string shown before the count of inserted lines.
	Deletions  string // Deletions is the string shown before the count of deleted lines.
}

type styles struct {
	Clear string // Clear is the style string that clears all styles.

	State  string // State is the style string printed before eventual special state.
	Branch string // Branch is the style string printed before the local branch.
	Remote string // Remote is the style string printed before the upstream branch.

	Divergence string // Divergence is the style string printed before divergence count/symbols.

	Staged    string // Staged is the style string printed before the staged files count.
	Conflict  string // Conflict is the style string printed before the conflict count.
	Modified  string // Modified is the style string printed before the modified files count.
	Untracked string // Untracked is the style string printed before the untracked files count.
	Stashed   string // Stashed is the style string printed before the stash entries count.
	Clean     string // Clean is the style string printed before the clean symbols.

	Insertions string // Insertions is the style string printed before the count of inserted lines.
	Deletions  string // Deletions is the style string printed before the count of deleted lines.
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
	return nil
}

type options struct {
	BranchMaxLen int       `yaml:"branch_max_len"`
	BranchTrim   direction `yaml:"branch_trim"`
	Ellipsis     string    `yaml:"ellipsis"`
	HideClean    bool      `yaml:"hide_clean"`
}

// A Formater formats git status to a tmux style string.
type Formater struct {
	Config
	b  bytes.Buffer
	st *gitstatus.Status
}

// truncate returns s, truncated so that it is no more than max runes long.
// Depending on the provided direction, truncation is performed right or left.
// If s is returned truncated, the truncated part is replaced with the
// 'ellipsis' string.
//
// If max is zero, negative or greater than the number of runes in s, truncate
// just returns s.
//
// NOTE: If max is lower than len(ellipsis), in other words it we're not even
// allowed to just return the ellipsis string, then we just return the maximum
// number of runes we can, without inserting ellpisis.
func truncate(s, ellipsis string, max int, dir direction) string {
	slen := utf8.RuneCountInString(s)
	if max <= 0 || slen <= max {
		return s
	}

	runes := []rune(s)
	ell := []rune(ellipsis)

	if max < len(ellipsis) {
		ell = nil // Just truncate s since even ellipsis don't fit.
	}

	switch dir {
	case dirRight:
		runes = runes[:max-len(ell)]
		runes = append(runes, ell...)
	case dirLeft:
		runes = runes[len(runes)+len(ell)-max:]
		runes = append(ell, runes...)
	}
	return string(runes)
}

// Format writes st as json into w.
func (f *Formater) Format(w io.Writer, st *gitstatus.Status) error {
	f.st = st
	f.clear()

	// overall working tree state
	if f.st.IsInitial {
		branch := truncate(f.st.LocalBranch, f.Options.Ellipsis, f.Options.BranchMaxLen, f.Options.BranchTrim)
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
			f.remoteBranch()
			f.divergence()
		case "remote-branch":
			f.remoteBranch()
		case "divergence":
			f.divergence()
		case "flags":
			f.flags()
		case "stats":
			f.stats()
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

func (f *Formater) remoteBranch() {
	if f.st.RemoteBranch == "" {
		return
	}

	f.clear()

	branch := truncate(f.st.RemoteBranch, f.Options.Ellipsis, f.Options.BranchMaxLen, f.Options.BranchTrim)
	fmt.Fprintf(&f.b, "%s%s", f.Styles.Remote, branch)
	f.b.WriteString(" ")
}

func (f *Formater) divergence() {
	if f.st.BehindCount == 0 && f.st.AheadCount == 0 {
		return
	}

	f.clear()
	fmt.Fprintf(&f.b, "%s", f.Styles.Divergence)

	if f.st.BehindCount != 0 {
		fmt.Fprintf(&f.b, "%s%d", f.Symbols.Behind, f.st.BehindCount)
	}

	if f.st.AheadCount != 0 {
		fmt.Fprintf(&f.b, "%s%d", f.Symbols.Ahead, f.st.AheadCount)
	}

	if f.st.BehindCount != 0 || f.st.AheadCount != 0 {
		f.b.WriteString(" ")
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
		f.b.WriteString(" ")
		return
	}

	branch := truncate(f.st.LocalBranch, f.Options.Ellipsis, f.Options.BranchMaxLen, f.Options.BranchTrim)
	fmt.Fprintf(&f.b, "%s%s", f.Styles.Branch, branch)
	f.b.WriteString(" ")
}

func (f *Formater) flags() {
	var flags []string
	if f.st.IsClean {
		if f.st.NumStashed != 0 {
			flags = append(flags,
				fmt.Sprintf("%s%s%d", f.Styles.Stashed, f.Symbols.Stashed, f.st.NumStashed))
		}

		if f.Options.HideClean != true {
			flags = append(flags, fmt.Sprintf("%s%s", f.Styles.Clean, f.Symbols.Clean))
		}

		f.clear()
		f.b.WriteString(strings.Join(flags, " "))
		return
	}

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
		f.b.WriteString(" ")
	}
}

func (f *Formater) stats() {
	stats := make([]string, 0, 2)

	if f.st.Insertions != 0 {
		stats = append(stats, fmt.Sprintf("%s%s%d", f.Styles.Insertions, f.Symbols.Insertions, f.st.Insertions))
	}

	if f.st.Deletions != 0 {
		stats = append(stats, fmt.Sprintf("%s%s%d", f.Styles.Deletions, f.Symbols.Deletions, f.st.Deletions))
	}

	if len(stats) != 0 {
		f.clear()
		f.b.WriteString(strings.Join(stats, " "))
	}
}

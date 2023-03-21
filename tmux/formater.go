package tmux

import (
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
	// overall working tree state
	if f.st.IsInitial {
		branch := truncate(f.st.LocalBranch, f.Options.Ellipsis, f.Options.BranchMaxLen, f.Options.BranchTrim)
		s := fmt.Sprintf("%s%s%s [no commits yet]", f.Styles.Clear, f.Styles.Branch, branch)
		if flags := f.flags(); flags != "" {
			s += flags
		}
		_, err := io.WriteString(w, s)
		return err
	}

	_, err := fmt.Fprintf(w, "%s%s", f.Styles.Clear, f.format())
	return err
}

func (f *Formater) format() string {
	var ss []string
	for _, item := range f.Layout {
		switch item {
		case "branch":
			ss = append(ss, f.specialState())
		case "remote":
			ss = append(ss, f.remoteBranch(), f.divergence())
		case "remote-branch":
			ss = append(ss, f.remoteBranch())
		case "divergence":
			ss = append(ss, f.divergence())
		case "flags":
			ss = append(ss, f.flags())
		case "stats":
			ss = append(ss, f.stats())
		default:
			ss = append(ss, f.Styles.Clear+item)
		}
	}

	// Remove empty status elements.
	i := 0
	for _, s := range ss {
		if s != "" {
			ss[i] = s
			i++
		}
	}
	ss = ss[:i]

	// TODO(arl) replace with stings.Join(ss, " ") after having removed all " " added here and there
	return strings.Join(ss, "")
}

func (f *Formater) specialState() string {
	s := f.Styles.Clear

	switch f.st.State {
	case gitstatus.Rebasing:
		s += fmt.Sprintf("%s[rebase] ", f.Styles.State)
	case gitstatus.AM:
		s += fmt.Sprintf("%s[am] ", f.Styles.State)
	case gitstatus.AMRebase:
		s += fmt.Sprintf("%s[am-rebase] ", f.Styles.State)
	case gitstatus.Merging:
		s += fmt.Sprintf("%s[merge] ", f.Styles.State)
	case gitstatus.CherryPicking:
		s += fmt.Sprintf("%s[cherry-pick] ", f.Styles.State)
	case gitstatus.Reverting:
		s += fmt.Sprintf("%s[revert] ", f.Styles.State)
	case gitstatus.Bisecting:
		s += fmt.Sprintf("%s[bisect] ", f.Styles.State)
	case gitstatus.Default:
		s += fmt.Sprintf("%s%s", f.Styles.Branch, f.Symbols.Branch)
	}

	s += f.currentRef()
	return s
}

func (f *Formater) remoteBranch() string {
	if f.st.RemoteBranch == "" {
		return ""
	}

	s := f.Styles.Clear

	branch := truncate(f.st.RemoteBranch, f.Options.Ellipsis, f.Options.BranchMaxLen, f.Options.BranchTrim)
	s += fmt.Sprintf("%s%s", f.Styles.Remote, branch)
	s += " "
	return s
}

func (f *Formater) divergence() string {
	if f.st.BehindCount == 0 && f.st.AheadCount == 0 {
		return ""
	}

	s := f.Styles.Clear
	s += fmt.Sprintf("%s", f.Styles.Divergence)

	if f.st.BehindCount != 0 {
		s += fmt.Sprintf("%s%d", f.Symbols.Behind, f.st.BehindCount)
	}

	if f.st.AheadCount != 0 {
		s += fmt.Sprintf("%s%d", f.Symbols.Ahead, f.st.AheadCount)
	}

	if f.st.BehindCount != 0 || f.st.AheadCount != 0 {
		s += " "
	}
	return s
}

// func (f *Formater) clear() {
// 	// clear global style
// 	f.b.WriteString(f.Styles.Clear)
// }

func (f *Formater) currentRef() string {
	s := f.Styles.Clear

	if f.st.IsDetached {
		s += fmt.Sprintf("%s%s%s", f.Styles.Branch, f.Symbols.HashPrefix, f.st.HEAD)
		s += " "
		return s
	}

	branch := truncate(f.st.LocalBranch, f.Options.Ellipsis, f.Options.BranchMaxLen, f.Options.BranchTrim)
	s += fmt.Sprintf("%s%s", f.Styles.Branch, branch)
	s += " "
	return s
}

func (f *Formater) flags() string {
	var flags []string
	if f.st.IsClean {
		if f.st.NumStashed != 0 {
			flags = append(flags,
				fmt.Sprintf("%s%s%d", f.Styles.Stashed, f.Symbols.Stashed, f.st.NumStashed))
		}

		if !f.Options.HideClean {
			flags = append(flags, fmt.Sprintf("%s%s", f.Styles.Clean, f.Symbols.Clean))
		}
		s := f.Styles.Clear + strings.Join(flags, " ")
		s += " "
		return s
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
		s := f.Styles.Clear + strings.Join(flags, " ")
		s += " "
		return s
	}

	return ""
}

func (f *Formater) stats() string {
	stats := make([]string, 0, 2)

	if f.st.Insertions != 0 {
		stats = append(stats, fmt.Sprintf("%s%s%d", f.Styles.Insertions, f.Symbols.Insertions, f.st.Insertions))
	}

	if f.st.Deletions != 0 {
		stats = append(stats, fmt.Sprintf("%s%s%d", f.Styles.Deletions, f.Symbols.Deletions, f.st.Deletions))
	}

	if len(stats) != 0 {
		s := f.Styles.Clear + strings.Join(stats, " ")
		s += " "
		return s
	}
	return ""
}

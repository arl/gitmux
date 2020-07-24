package tmux

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arl/gitstatus"
)

func TestFormater_flags(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		layout  []string
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "clean flag",
			styles: styles{
				Clean: "StyleClean",
			},
			symbols: symbols{
				Clean: "SymbolClean",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			want: clear + "StyleCleanSymbolClean",
		},
		{
			name: "mixed flags",
			styles: styles{
				Modified: "StyleMod",
				Stashed:  "StyleStash",
				Staged:   "StyleStaged",
			},
			symbols: symbols{
				Modified: "SymbolMod",
				Stashed:  "SymbolStash",
				Staged:   "SymbolStaged",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
					NumStaged:   3,
				},
			},
			want: clear + "StyleStagedSymbolStaged3 StyleModSymbolMod2 StyleStashSymbolStash1",
		},
		{
			name: "mixed flags 2",
			styles: styles{
				Conflict:  "StyleConflict",
				Untracked: "StyleUntracked",
			},
			symbols: symbols{
				Conflict:  "SymbolConflict",
				Untracked: "SymbolUntracked",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumConflicts: 42,
					NumUntracked: 17,
				},
			},
			want: clear + "StyleConflictSymbolConflict42 StyleUntrackedSymbolUntracked17",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tc.styles, Symbols: tc.symbols, Layout: tc.layout},
				st:     tc.st,
			}
			f.flags()
			require.EqualValues(t, tc.want, f.b.String())
		})
	}
}

func TestFormater_divergence(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "no divergence",
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  0,
					BehindCount: 0,
				},
			},
			want: clear,
		},
		{
			name: "ahead only",
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  4,
					BehindCount: 0,
				},
			},
			want: clear + " ↓·4",
		},
		{
			name: "behind only",
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  0,
					BehindCount: 12,
				},
			},
			want: clear + " ↑·12",
		},
		{
			name: "diverged both ways",
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  41,
					BehindCount: 128,
				},
			},
			want: clear + " ↑·128↓·41",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tc.styles, Symbols: tc.symbols},
				st:     tc.st,
			}
			f.divergence()
			require.EqualValues(t, tc.want, f.b.String())
		})
	}
}

func TestFormater_BranchMaxLen(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		options options
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "no limit",
			options: options{
				BranchMaxLen: 0,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "foo/bar-baz",
					RemoteBranch: "remote/foo/bar-baz",
				},
			},
			want: clear + "foo/bar-baz" + clear + "remote/foo/bar-baz",
		},
		{
			name: "no truncate",
			options: options{
				BranchMaxLen: 11,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "foo/bar-baz",
					RemoteBranch: "remote/foo/bar-baz",
				},
			},
			want: clear + "foo/bar-baz" + clear + "remote/foo/bar-baz",
		},
		{
			name: "truncate",
			options: options{
				BranchMaxLen: 10,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "foo/bar-baz",
					RemoteBranch: "remote/foo/bar-baz",
				},
			},
			want: clear + "foo/bar..." + clear + "remote/foo/bar...",
		},
		{
			name: "truncate to 1",
			options: options{
				BranchMaxLen: 1,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "foo/bar-baz",
					RemoteBranch: "remote/foo/bar-baz",
				},
			},
			want: clear + "." + clear + "remote/.",
		},
		{
			name: "truncate utf-8 name",
			options: options{
				BranchMaxLen: 9,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "foo/测试这个名字",
					RemoteBranch: "remote/foo/测试这个名字",
				},
			},
			want: clear + "foo/测试..." + clear + "remote/foo/测试...",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Options: tc.options},
				st:     tc.st,
			}
			f.currentRef()
			f.remoteBranch()
			require.EqualValues(t, tc.want, f.b.String())
		})
	}
}

func TestFormater_Format(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		layout  []string
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "default format",
			styles: styles{
				Clean:    "StyleClean",
				Branch:   "StyleBranch",
				Modified: "StyleMod",
				Remote:   "StyleRemote",
			},
			symbols: symbols{
				Branch:   "SymbolBranch",
				Clean:    "SymbolClean",
				Modified: "SymbolMod",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					NumModified:  2,
				},
			},
			want: clear + "StyleBranchSymbolBranch" + clear + "Local" + ".." + clear + "StyleRemoteRemote" + clear + " - " + clear + "StyleModSymbolMod2",
		},
		{
			name: "branch, different delimiter, flags",
			styles: styles{
				Branch:   "StyleBranch",
				Remote:   "StyleRemote",
				Modified: "StyleMod",
			},
			symbols: symbols{
				Branch:   "SymbolBranch",
				Ahead:    "SymbolAhead",
				Modified: "SymbolMod",
			},
			layout: []string{"branch", " ~~ ", "flags"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					NumModified:  2,
					AheadCount:   1,
				},
			},
			want: clear + "StyleBranchSymbolBranch" + clear + "Local" + " ~~ " + clear + "StyleModSymbolMod2",
		},
		{
			name: "remote only",
			styles: styles{
				Branch: "StyleBranch",
				Remote: "StyleRemote",
			},
			symbols: symbols{
				Branch: "SymbolBranch",
				Ahead:  "SymbolAhead",
			},
			layout: []string{"remote"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					AheadCount:   1,
				},
			},
			want: clear + "StyleRemoteRemote" + clear + " SymbolAhead1",
		},
		{
			name: "empty",
			styles: styles{
				Branch:   "StyleBranch",
				Modified: "StyleMod",
			},
			symbols: symbols{
				Branch:   "SymbolBranch",
				Modified: "SymbolMod",
			},
			layout: []string{},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "Local",
					NumModified: 2,
				},
			},
			want: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tc.styles, Symbols: tc.symbols, Layout: tc.layout},
			}

			f.Format(os.Stdout, tc.st)
			f.format()
			require.EqualValues(t, tc.want, f.b.String())
		})
	}
}

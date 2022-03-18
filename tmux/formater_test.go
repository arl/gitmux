package tmux

import (
	"io"
	"testing"

	"github.com/arl/gitstatus"
)

func TestFlags(t *testing.T) {
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
				Clear: "StyleClear",
				Clean: "StyleClean",
			},
			symbols: symbols{
				Clean: "SymbolClean",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			want: "StyleClear" + "StyleCleanSymbolClean",
		},
		{
			name: "stash + clean flag",
			styles: styles{
				Clear:   "StyleClear",
				Clean:   "StyleClean",
				Stashed: "StyleStash",
			},
			symbols: symbols{
				Clean:   "SymbolClean",
				Stashed: "SymbolStash",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				IsClean:    true,
				NumStashed: 1,
			},
			want: "StyleClearStyleStashSymbolStash1 StyleCleanSymbolClean",
		},
		{
			name: "mixed flags",
			styles: styles{
				Clear:    "StyleClear",
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
			want: "StyleClear" + "StyleStagedSymbolStaged3 StyleModSymbolMod2 StyleStashSymbolStash1",
		},
		{
			name: "mixed flags 2",
			styles: styles{
				Clear:     "StyleClear",
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
			want: "StyleClear" + "StyleConflictSymbolConflict42 StyleUntrackedSymbolUntracked17",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Layout: tt.layout},
				st:     tt.st,
			}
			f.flags()

			if got := f.b.String(); got != tt.want {
				t.Errorf("got:\n%s\n\nwant:\n%s\n", got, tt.want)
			}
		})
	}
}

func TestDivergence(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "no divergence",
			styles: styles{
				Clear: "StyleClear",
			},
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
			want: "",
		},
		{
			name: "ahead only",
			styles: styles{
				Clear: "StyleClear",
			},
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
			want: "StyleClear" + " ↓·4",
		},
		{
			name: "behind only",
			styles: styles{
				Clear: "StyleClear",
			},
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
			want: "StyleClear" + " ↑·12",
		},
		{
			name: "diverged both ways",
			styles: styles{
				Clear: "StyleClear",
			},
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
			want: "StyleClear" + " ↑·128↓·41",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols},
				st:     tt.st,
			}
			f.divergence()

			if got := f.b.String(); got != tt.want {
				t.Errorf("got:\n%s\n\nwant:\n%s\n", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s    string
		max  int
		dir  direction
		want string
	}{
		/* trim right */
		{
			s:    "br",
			max:  1,
			dir:  dirRight,
			want: "b",
		},
		{
			s:    "br",
			max:  3,
			dir:  dirRight,
			want: "br",
		},
		{
			s:    "super-long-branch",
			max:  3,
			dir:  dirRight,
			want: "...",
		},
		{
			s:    "super-long-branch",
			max:  15,
			dir:  dirRight,
			want: "super-long-b...",
		},
		{
			s:    "super-long-branch",
			max:  17,
			dir:  dirRight,
			want: "super-long-branch",
		},
		{
			s:    "长長的-树樹枝",
			max:  6,
			dir:  dirRight,
			want: "长長的...",
		},
		{
			s:    "super-long-branch",
			max:  32,
			dir:  dirRight,
			want: "super-long-branch",
		},
		{
			s:    "super-long-branch",
			max:  0,
			dir:  dirRight,
			want: "super-long-branch",
		},
		{
			s:    "super-long-branch",
			max:  -1,
			dir:  dirRight,
			want: "super-long-branch",
		},

		/* trim left */
		{
			s:    "br",
			max:  1,
			dir:  dirLeft,
			want: "r",
		},
		{
			s:    "br",
			max:  3,
			dir:  dirLeft,
			want: "br",
		},
		{
			s:    "super-long-branch",
			max:  3,
			dir:  dirLeft,
			want: "...",
		},
		{
			s:    "super-long-branch",
			max:  15,
			dir:  dirLeft,
			want: "...-long-branch",
		},
		{
			s:    "super-long-branch",
			max:  17,
			dir:  dirLeft,
			want: "super-long-branch",
		},
		{
			s:    "长長的-树樹枝",
			max:  6,
			dir:  dirLeft,
			want: "...树樹枝",
		},
		{
			s:    "super-long-branch",
			max:  32,
			dir:  dirLeft,
			want: "super-long-branch",
		},
		{
			s:    "super-long-branch",
			max:  0,
			dir:  dirLeft,
			want: "super-long-branch",
		},
		{
			s:    "super-long-branch",
			max:  -1,
			dir:  dirLeft,
			want: "super-long-branch",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := truncate(tt.s, tt.max, tt.dir); got != tt.want {
				t.Errorf("truncate(%q, %d, %s) = %q, want %q", tt.s, tt.max, tt.dir, got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		layout  []string
		options options
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "default format",
			styles: styles{
				Clear:    "StyleClear",
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
			want: "StyleClear" + "StyleBranchSymbolBranch" +
				"StyleClear" + "StyleBranch" + "Local" +
				"StyleClear" + ".." +
				"StyleClear" + "StyleRemoteRemote" +
				"StyleClear" + " - " +
				"StyleClear" + "StyleModSymbolMod2",
		},
		{
			name: "branch, different delimiter, flags",
			styles: styles{
				Clear:    "StyleClear",
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
			want: "StyleClear" + "StyleBranchSymbolBranch" +
				"StyleClear" + "StyleBranch" + "Local" +
				"StyleClear" + " ~~ " +
				"StyleClear" + "StyleModSymbolMod2",
		},
		{
			name: "remote only",
			styles: styles{
				Clear:  "StyleClear",
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
			want: "StyleClear" + "StyleRemoteRemote" +
				"StyleClear" + " SymbolAhead1",
		},
		{
			name: "empty",
			styles: styles{
				Clear:    "StyleClear",
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
		{
			name: "branch and remote, branch_max_len not zero",
			styles: styles{
				Clear:  "StyleClear",
				Branch: "StyleBranch",
				Remote: "StyleRemote",
			},
			symbols: symbols{
				Branch: "SymbolBranch",
			},
			layout: []string{"branch", " ", "remote"},
			options: options{
				BranchMaxLen: 9,
				BranchTrim:   dirRight,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "branchName",
					RemoteBranch: "remote/branchName",
				},
			},
			want: "StyleClear" + "StyleBranch" + "SymbolBranch" +
				"StyleClear" + "StyleBranch" + "branch..." +
				"StyleClear" + " " +
				"StyleClear" + "StyleRemote" + "remote...",
		},
		{
			name: "branch and remote, branch_max_len not zero and trim left",
			styles: styles{
				Clear:  "StyleClear",
				Branch: "StyleBranch",
				Remote: "StyleRemote",
			},
			symbols: symbols{
				Branch: "SymbolBranch",
			},
			layout: []string{"branch", " ", "remote"},
			options: options{
				BranchMaxLen: 9,
				BranchTrim:   dirLeft,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "nameBranch",
					RemoteBranch: "remote/nameBranch",
				},
			},
			want: "StyleClear" + "StyleBranch" + "SymbolBranch" +
				"StyleClear" + "StyleBranch" + "...Branch" +
				"StyleClear" + " " +
				"StyleClear" + "StyleRemote" + "...Branch",
		},
		{
			name: "issue-32",
			styles: styles{
				Clear:  "StyleClear",
				Branch: "StyleBranch",
			},
			symbols: symbols{
				Branch: "SymbolBranch",
			},
			layout: []string{"branch"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "branchName",
				},
			},
			want: "StyleClear" + "StyleBranch" + "SymbolBranch" +
				"StyleClear" + "StyleBranch" + "branchName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Layout: tt.layout, Options: tt.options},
			}

			if err := f.Format(io.Discard, tt.st); err != nil {
				t.Fatalf("Format error: %s", err)
			}

			f.format()
			if got := f.b.String(); got != tt.want {
				t.Errorf("got:\n%s\n\nwant:\n%s\n", got, tt.want)
			}
		})
	}
}

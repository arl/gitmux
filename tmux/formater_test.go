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
			layout: []string{"branch", "..", "remote", "- ", "flags"},
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
			layout: []string{"branch", "..", "remote", "- ", "flags"},
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

			compareStrings(t, tt.want, f.flags())
		})
	}
}

func TestDivergence(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		options options
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "no divergence",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
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
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
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
			want: "StyleClearStyleDivergence" + "↓·4",
		},
		{
			name: "behind only",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
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
			want: "StyleClearStyleDivergence" + "↑·12",
		},
		{
			name: "diverged both ways",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
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
			want: "StyleClearStyleDivergence" + "↑·128↓·41",
		},
		{
			name: "divergence-space:true and ahead:0",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
			},
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			options: options{
				DivergenceSpace: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  0,
					BehindCount: 12,
				},
			},
			want: "StyleClearStyleDivergence" + "↑·12",
		},
		{
			name: "divergence-space:false and diverged both ways",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
			},
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			options: options{
				DivergenceSpace: true,
				SwapDivergence:  false,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  41,
					BehindCount: 128,
				},
			},
			want: "StyleClearStyleDivergence" + "↑·128 ↓·41",
		},
		{
			name: "divergence-space:true and diverged both ways",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
			},
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			options: options{
				DivergenceSpace: true,
				SwapDivergence:  true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  41,
					BehindCount: 128,
				},
			},
			want: "StyleClearStyleDivergence" + "↓·41 ↑·128",
		},
		{
			name: "swap divergence ahead only",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
			},
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			options: options{
				SwapDivergence: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  4,
					BehindCount: 0,
				},
			},
			want: "StyleClearStyleDivergence" + "↓·4",
		},
		{
			name: "swap divergence behind only",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
			},
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			options: options{
				SwapDivergence: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  0,
					BehindCount: 12,
				},
			},
			want: "StyleClearStyleDivergence" + "↑·12",
		},
		{
			name: "swap divergence both ways",
			styles: styles{
				Clear:      "StyleClear",
				Divergence: "StyleDivergence",
			},
			symbols: symbols{
				Ahead:  "↓·",
				Behind: "↑·",
			},
			options: options{
				SwapDivergence: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					AheadCount:  41,
					BehindCount: 128,
				},
			},
			want: "StyleClearStyleDivergence" + "↓·41↑·128",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Options: tt.options},
				st:     tt.st,
			}

			compareStrings(t, tt.want, f.divergence())
		})
	}
}

func Test_truncate(t *testing.T) {
	tests := []struct {
		s        string
		max      int
		ellipsis string
		dir      direction
		want     string
	}{
		/* trim right */
		{
			s:        "br",
			ellipsis: "...",
			max:      1,
			dir:      dirRight,
			want:     "b",
		},
		{
			s:        "br",
			ellipsis: "...",
			max:      3,
			dir:      dirRight,
			want:     "br",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      3,
			dir:      dirRight,
			want:     "...",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      15,
			dir:      dirRight,
			want:     "super-long-b...",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      17,
			dir:      dirRight,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "…",
			max:      17,
			dir:      dirRight,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "…",
			max:      15,
			dir:      dirRight,
			want:     "super-long-bra…",
		},
		{
			s:        "长長的-树樹枝",
			ellipsis: "...",
			max:      6,
			dir:      dirRight,
			want:     "长長的...",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      32,
			dir:      dirRight,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      0,
			dir:      dirRight,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      -1,
			dir:      dirRight,
			want:     "super-long-branch",
		},

		/* trim left */
		{
			s:        "br",
			ellipsis: "...",
			max:      1,
			dir:      dirLeft,
			want:     "r",
		},
		{
			s:        "br",
			ellipsis: "",
			max:      1,
			dir:      dirLeft,
			want:     "r",
		},
		{
			s:        "br",
			ellipsis: "...",
			max:      3,
			dir:      dirLeft,
			want:     "br",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      3,
			dir:      dirLeft,
			want:     "...",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      15,
			dir:      dirLeft,
			want:     "...-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      17,
			dir:      dirLeft,
			want:     "super-long-branch",
		},
		{
			s:        "长長的-树樹枝",
			ellipsis: "...",
			max:      6,
			dir:      dirLeft,
			want:     "...树樹枝",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      32,
			dir:      dirLeft,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      0,
			dir:      dirLeft,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      -1,
			dir:      dirLeft,
			want:     "super-long-branch",
		},

		/* trim center */
		{
			s:        "br",
			ellipsis: "...",
			max:      1,
			dir:      dirCenter,
			want:     "r",
		},
		{
			s:        "br",
			ellipsis: "",
			max:      1,
			dir:      dirCenter,
			want:     "r",
		},
		{
			s:        "br",
			ellipsis: "...",
			max:      3,
			dir:      dirCenter,
			want:     "br",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      3,
			dir:      dirCenter,
			want:     "...",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      15,
			dir:      dirCenter,
			want:     "super-...branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      17,
			dir:      dirCenter,
			want:     "super-long-branch",
		},
		{
			s:        "长長的-树樹枝",
			ellipsis: "...",
			max:      6,
			dir:      dirCenter,
			want:     "长...樹枝",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      32,
			dir:      dirCenter,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      0,
			dir:      dirCenter,
			want:     "super-long-branch",
		},
		{
			s:        "super-long-branch",
			ellipsis: "...",
			max:      -1,
			dir:      dirCenter,
			want:     "super-long-branch",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			compareStrings(t, tt.want, truncate(tt.s, tt.ellipsis, tt.max, tt.dir))
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
			layout: []string{"branch", " .. ", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					NumModified:  2,
				},
			},
			want: "StyleClear" + "StyleBranchSymbolBranch" +
				"StyleClear" + "StyleBranch" + "Local" +
				"StyleClear" + " .. " +
				"StyleClear" + "StyleRemoteRemote" +
				"StyleClear" + " - " +
				"StyleClear" + "StyleModSymbolMod2" +
				resetStyles,
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
			layout: []string{"branch", "~~", "flags"},
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
				"StyleClear" + "~~" +
				"StyleClear" + "StyleModSymbolMod2" +
				resetStyles,
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
			want: "StyleClear" + "StyleRemoteRemote " +
				"StyleClear" + "SymbolAhead1" +
				resetStyles,
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
			want: resetStyles,
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
			layout: []string{"branch", "/", "remote"},
			options: options{
				BranchMaxLen: 9,
				BranchTrim:   dirRight,
				Ellipsis:     `…`,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "branchName",
					RemoteBranch: "remote/branchName",
				},
			},
			want: "StyleClear" + "StyleBranch" + "SymbolBranch" +
				"StyleClear" + "StyleBranch" + "branchNa…" +
				"StyleClear" + "/" +
				"StyleClear" + "StyleRemote" + "remote/b…" +
				resetStyles,
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
			layout: []string{"branch", "remote"},
			options: options{
				BranchMaxLen: 9,
				BranchTrim:   dirLeft,
				Ellipsis:     "...",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "nameBranch",
					RemoteBranch: "remote/nameBranch",
				},
			},
			want: "StyleClear" + "StyleBranch" + "SymbolBranch" +
				"StyleClear" + "StyleBranch" + "...Branch " +
				"StyleClear" + "StyleRemote" + "...Branch" +
				resetStyles,
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
				"StyleClear" + "StyleBranch" + "branchName" +
				resetStyles,
		},
		{
			name: "hide clean option true",
			styles: styles{
				Clear: "StyleClear",
				Clean: "StyleClean",
			},
			symbols: symbols{
				Clean: "SymbolClean",
			},
			layout: []string{"flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			options: options{
				HideClean: true,
			},
			want: resetStyles,
		},
		{
			name: "hide clean option false",
			styles: styles{
				Clear: "StyleClear",
				Clean: "StyleClean",
			},
			symbols: symbols{
				Clean: "SymbolClean",
			},
			layout: []string{"flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			options: options{
				HideClean: false,
			},
			want: "StyleClear" + "StyleCleanSymbolClean" + resetStyles,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Layout: tt.layout, Options: tt.options},
			}

			if err := f.Format(io.Discard, tt.st); err != nil {
				t.Fatalf("Format error: %s", err)
				return
			}

			compareStrings(t, tt.want, f.format())
		})
	}
}

func Test_stats(t *testing.T) {
	tests := []struct {
		name                  string
		layout                []string
		insertions, deletions int
		want                  string
	}{
		{
			name: "nothing",
			want: "",
		},
		{
			name:       "insertions",
			insertions: 12,
			want:       "StyleClear" + "StyleInsertionsSymbolInsertions12",
		},
		{
			name:      "deletions",
			deletions: 12,
			want:      "StyleClear" + "StyleDeletionsSymbolDeletions12",
		},
		{
			name:       "insertions and deletions",
			insertions: 1,
			deletions:  2,
			want:       "StyleClear" + "StyleInsertionsSymbolInsertions1" + " " + "StyleDeletionsSymbolDeletions2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{
					Styles: styles{
						Clear:      "StyleClear",
						Deletions:  "StyleDeletions",
						Insertions: "StyleInsertions",
					},
					Symbols: symbols{
						Deletions:  "SymbolDeletions",
						Insertions: "SymbolInsertions",
					},
					Layout: []string{"stats"},
				},
				st: &gitstatus.Status{
					Insertions: tt.insertions,
					Deletions:  tt.deletions,
				},
			}

			compareStrings(t, tt.want, f.stats())
		})
	}
}

func compareStrings(t *testing.T, want, got string) {
	if got != want {
		t.Errorf(`
	got:
%q

	want:
%q`, got, want)
	}
}

package tmux

import (
	"testing"

	"github.com/arl/gitstatus"
	"github.com/stretchr/testify/require"
)

func TestFormater_flags(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "clean flag",
			styles: styles{
				Clean: "CleanStyle",
			},
			symbols: symbols{
				Clean:     "CleanSymbol",
				Delimiter: " - ",
			},
			st: &gitstatus.Status{
				IsClean: true,
			},
			want: clear + " - CleanStyleCleanSymbol",
		},
		{
			name: "mixed flags",
			styles: styles{
				Modified: "StyleMod",
				Stashed:  "StyleStash",
				Staged:   "StyleStaged",
			},
			symbols: symbols{
				Modified:  "SymbolMod",
				Stashed:   "SymbolStash",
				Staged:    "SymbolStaged",
				Delimiter: " - ",
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
					NumStaged:   3,
				},
			},
			want: clear + " - StyleStagedSymbolStaged3 StyleModSymbolMod2 StyleStashSymbolStash1",
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
				Delimiter: " - ",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumConflicts: 42,
					NumUntracked: 17,
				},
			},
			want: clear + " - StyleConflictSymbolConflict42 StyleUntrackedSymbolUntracked17",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tc.styles, Symbols: tc.symbols},
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

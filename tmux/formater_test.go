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
				Clear: "[style:clear]",
				Clean: "[style:clean]",
			},
			symbols: symbols{
				Clean: "[symbol:clean]",
			},
			layout: []string{"branch", "..", "remote", "- ", "flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			want: "[style:clear]" + "[style:clean][symbol:clean]",
		},
		{
			name: "stash + clean flag",
			styles: styles{
				Clear:   "[style:clear]",
				Clean:   "[style:clean]",
				Stashed: "[style:stashed]",
			},
			symbols: symbols{
				Clean:   "[symbol:clean]",
				Stashed: "[symbol:stashed]",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				IsClean:    true,
				NumStashed: 1,
			},
			want: "[style:clear][style:stashed][symbol:stashed]1 [style:clean][symbol:clean]",
		},
		{
			name: "mixed flags",
			styles: styles{
				Clear:    "[style:clear]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
				Staged:   "[style:staged]",
			},
			symbols: symbols{
				Modified: "[symbol:mod]",
				Stashed:  "[symbol:stashed]",
				Staged:   "[symbol:staged]",
			},
			layout: []string{"branch", "..", "remote", "- ", "flags"},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
					NumStaged:   3,
				},
			},
			want: "[style:clear]" + "[style:staged][symbol:staged]3 [style:mod][symbol:mod]2 [style:stashed][symbol:stashed]1",
		},
		{
			name: "mixed flags 2",
			styles: styles{
				Clear:     "[style:clear]",
				Conflict:  "[style:conflict]",
				Untracked: "[style:untracked]",
			},
			symbols: symbols{
				Conflict:  "[symbol:conflict]",
				Untracked: "[symbol:untracked]",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumConflicts: 42,
					NumUntracked: 17,
				},
			},
			want: "[style:clear]" + "[style:conflict][symbol:conflict]42 [style:untracked][symbol:untracked]17",
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

func TestFlagsWithoutCountBehavior(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		options options
		st      *gitstatus.Status
		want    string
	}{
		// Case 1: non-empty symbol, count=1, flags_without_count=false
		{
			name: "case 1a: non-empty symbol, count=1, flags_without_count=false",
			styles: styles{
				Clear:  "[style:clear]",
				Staged: "[style:staged]",
			},
			symbols: symbols{
				Staged: "[symbol:staged]",
			},
			options: options{
				FlagsWithoutCount: false,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumStaged: 1,
				},
			},
			want: "[style:clear][style:staged][symbol:staged]1",
		},
		// Case 1: non-empty symbol, count=1, flags_without_count=true
		{
			name: "case 1b: non-empty symbol, count=1, flags_without_count=true",
			styles: styles{
				Clear:  "[style:clear]",
				Staged: "[style:staged]",
			},
			symbols: symbols{
				Staged: "[symbol:staged]",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumStaged: 1,
				},
			},
			want: "[style:clear][style:staged][symbol:staged]",
		},
		// Case 2: empty symbol, count=1, flags_without_count=false
		{
			name: "case 2a: empty symbol, count=1, flags_without_count=false",
			styles: styles{
				Clear:  "[style:clear]",
				Staged: "[style:staged]",
			},
			symbols: symbols{
				Staged: "",
			},
			options: options{
				FlagsWithoutCount: false,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumStaged: 1,
				},
			},
			want: "[style:clear][style:staged]1",
		},
		// Case 2: empty symbol, count=1, flags_without_count=true
		{
			name: "case 2b: empty symbol, count=1, flags_without_count=true",
			styles: styles{
				Clear:  "[style:clear]",
				Staged: "[style:staged]",
			},
			symbols: symbols{
				Staged: "",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumStaged: 1,
				},
			},
			want: "",
		},
		// Case 3: count=0, flags_without_count=false
		{
			name: "case 3a: count=0, flags_without_count=false",
			styles: styles{
				Clear:  "[style:clear]",
				Staged: "[style:staged]",
			},
			symbols: symbols{
				Staged: "[symbol:staged]",
			},
			options: options{
				FlagsWithoutCount: false,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumStaged: 0,
				},
			},
			want: "",
		},
		// Case 3: count=0, flags_without_count=true
		{
			name: "case 3b: count=0, flags_without_count=true",
			styles: styles{
				Clear:  "[style:clear]",
				Staged: "[style:staged]",
			},
			symbols: symbols{
				Staged: "[symbol:staged]",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumStaged: 0,
				},
			},
			want: "",
		},
		// Mixed case: multiple flags with different symbol states
		{
			name: "mixed flags: some empty symbols, some non-empty, flags_without_count=false",
			styles: styles{
				Clear:    "[style:clear]",
				Staged:   "[style:staged]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
			},
			symbols: symbols{
				Staged:   "",  // empty symbol, should show count only
				Modified: "M", // non-empty symbol, should show symbol+count
				Stashed:  "",  // empty symbol, should show count only
			},
			options: options{
				FlagsWithoutCount: false,
			},
			st: &gitstatus.Status{
				NumStashed: 2,
				Porcelain: gitstatus.Porcelain{
					NumStaged:   1,
					NumModified: 3,
				},
			},
			want: "[style:clear][style:staged]1 [style:mod]M3 [style:stashed]2",
		},
		// Mixed case: multiple flags with different symbol states, flags_without_count=true
		{
			name: "mixed flags: some empty symbols, some non-empty, flags_without_count=true",
			styles: styles{
				Clear:    "[style:clear]",
				Staged:   "[style:staged]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
			},
			symbols: symbols{
				Staged:   "",  // empty symbol, should show nothing
				Modified: "M", // non-empty symbol, should show symbol only
				Stashed:  "",  // empty symbol, should show nothing
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				NumStashed: 2,
				Porcelain: gitstatus.Porcelain{
					NumStaged:   1,
					NumModified: 3,
				},
			},
			want: "[style:clear][style:mod]M",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Options: tt.options},
				st:     tt.st,
			}

			compareStrings(t, tt.want, f.flags())
		})
	}
}

func TestFlagsWithoutCount(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		options options
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "flags with counts (default)",
			styles: styles{
				Clear:    "[style:clear]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
				Staged:   "[style:staged]",
			},
			symbols: symbols{
				Modified: "[symbol:mod]",
				Stashed:  "[symbol:stashed]",
				Staged:   "[symbol:staged]",
			},
			options: options{
				FlagsWithoutCount: false,
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
					NumStaged:   3,
				},
			},
			want: "[style:clear]" + "[style:staged][symbol:staged]3 [style:mod][symbol:mod]2 [style:stashed][symbol:stashed]1",
		},
		{
			name: "flags without counts",
			styles: styles{
				Clear:    "[style:clear]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
				Staged:   "[style:staged]",
			},
			symbols: symbols{
				Modified: "[symbol:mod]",
				Stashed:  "[symbol:stashed]",
				Staged:   "[symbol:staged]",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
					NumStaged:   3,
				},
			},
			want: "[style:clear]" + "[style:staged][symbol:staged] [style:mod][symbol:mod] [style:stashed][symbol:stashed]",
		},
		{
			name: "all flags without counts",
			styles: styles{
				Clear:     "[style:clear]",
				Conflict:  "[style:conflict]",
				Modified:  "[style:mod]",
				Stashed:   "[style:stashed]",
				Staged:    "[style:staged]",
				Untracked: "[style:untracked]",
			},
			symbols: symbols{
				Conflict:  "[symbol:conflict]",
				Modified:  "[symbol:mod]",
				Stashed:   "[symbol:stashed]",
				Staged:    "[symbol:staged]",
				Untracked: "[symbol:untracked]",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				NumStashed: 5,
				Porcelain: gitstatus.Porcelain{
					NumConflicts: 1,
					NumModified:  10,
					NumStaged:    3,
					NumUntracked: 7,
				},
			},
			want: "[style:clear]" + "[style:staged][symbol:staged] [style:conflict][symbol:conflict] [style:mod][symbol:mod] [style:stashed][symbol:stashed] [style:untracked][symbol:untracked]",
		},
		{
			name: "clean with stash without count",
			styles: styles{
				Clear:   "[style:clear]",
				Clean:   "[style:clean]",
				Stashed: "[style:stashed]",
			},
			symbols: symbols{
				Clean:   "[symbol:clean]",
				Stashed: "[symbol:stashed]",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				IsClean:    true,
				NumStashed: 1,
			},
			want: "[style:clear][style:stashed][symbol:stashed] [style:clean][symbol:clean]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Options: tt.options},
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
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↓·4",
		},
		{
			name: "behind only",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↑·12",
		},
		{
			name: "diverged both ways",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↑·128↓·41",
		},
		{
			name: "divergence-space:true and ahead:0",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↑·12",
		},
		{
			name: "divergence-space:false and diverged both ways",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↑·128 ↓·41",
		},
		{
			name: "divergence-space:true and diverged both ways",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↓·41 ↑·128",
		},
		{
			name: "swap divergence ahead only",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↓·4",
		},
		{
			name: "swap divergence behind only",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↑·12",
		},
		{
			name: "swap divergence both ways",
			styles: styles{
				Clear:      "[style:clear]",
				Divergence: "[style:divergence]",
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
			want: "[style:clear][style:divergence]" + "↓·41↑·128",
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
				Clear:    "[style:clear]",
				Clean:    "[style:clean]",
				Branch:   "[style:branch]",
				Modified: "[style:mod]",
				Remote:   "[style:remote]",
			},
			symbols: symbols{
				Branch:   "[symbol:branch]",
				Clean:    "[symbol:clean]",
				Modified: "[symbol:mod]",
			},
			layout: []string{"branch", " .. ", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					NumModified:  2,
				},
			},
			want: "[style:clear]" + "[style:branch][symbol:branch]" +
				"[style:clear]" + "[style:branch]" + "Local" +
				"[style:clear]" + " .. " +
				"[style:clear]" + "[style:remote]Remote" +
				"[style:clear]" + " - " +
				"[style:clear]" + "[style:mod][symbol:mod]2" +
				resetStyles,
		},
		{
			name: "branch, different delimiter, flags",
			styles: styles{
				Clear:    "[style:clear]",
				Branch:   "[style:branch]",
				Remote:   "[style:remote]",
				Modified: "[style:mod]",
			},
			symbols: symbols{
				Branch:   "[symbol:branch]",
				Ahead:    "[symbol:ahead]",
				Modified: "[symbol:mod]",
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
			want: "[style:clear]" + "[style:branch][symbol:branch]" +
				"[style:clear]" + "[style:branch]" + "Local" +
				"[style:clear]" + "~~" +
				"[style:clear]" + "[style:mod][symbol:mod]2" +
				resetStyles,
		},
		{
			name: "remote only",
			styles: styles{
				Clear:  "[style:clear]",
				Branch: "[style:branch]",
				Remote: "[style:remote]",
			},
			symbols: symbols{
				Branch: "[symbol:branch]",
				Ahead:  "[symbol:ahead]",
			},
			layout: []string{"remote"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					AheadCount:   1,
				},
			},
			want: "[style:clear]" + "[style:remote]Remote " +
				"[style:clear]" + "[symbol:ahead]1" +
				resetStyles,
		},
		{
			name: "empty",
			styles: styles{
				Clear:    "[style:clear]",
				Branch:   "[style:branch]",
				Modified: "[style:mod]",
			},
			symbols: symbols{
				Branch:   "[symbol:branch]",
				Modified: "[symbol:mod]",
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
				Clear:  "[style:clear]",
				Branch: "[style:branch]",
				Remote: "[style:remote]",
			},
			symbols: symbols{
				Branch: "[symbol:branch]",
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
			want: "[style:clear]" + "[style:branch]" + "[symbol:branch]" +
				"[style:clear]" + "[style:branch]" + "branchNa…" +
				"[style:clear]" + "/" +
				"[style:clear]" + "[style:remote]" + "remote/b…" +
				resetStyles,
		},
		{
			name: "branch and remote, branch_max_len not zero and trim left",
			styles: styles{
				Clear:  "[style:clear]",
				Branch: "[style:branch]",
				Remote: "[style:remote]",
			},
			symbols: symbols{
				Branch: "[symbol:branch]",
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
			want: "[style:clear]" + "[style:branch]" + "[symbol:branch]" +
				"[style:clear]" + "[style:branch]" + "...Branch " +
				"[style:clear]" + "[style:remote]" + "...Branch" +
				resetStyles,
		},
		{
			name: "issue-32",
			styles: styles{
				Clear:  "[style:clear]",
				Branch: "[style:branch]",
			},
			symbols: symbols{
				Branch: "[symbol:branch]",
			},
			layout: []string{"branch"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "branchName",
				},
			},
			want: "[style:clear]" + "[style:branch]" + "[symbol:branch]" +
				"[style:clear]" + "[style:branch]" + "branchName" +
				resetStyles,
		},
		{
			name: "hide clean option true",
			styles: styles{
				Clear: "[style:clear]",
				Clean: "[style:clean]",
			},
			symbols: symbols{
				Clean: "[symbol:clean]",
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
				Clear: "[style:clear]",
				Clean: "[style:clean]",
			},
			symbols: symbols{
				Clean: "[symbol:clean]",
			},
			layout: []string{"flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			options: options{
				HideClean: false,
			},
			want: "[style:clear]" + "[style:clean][symbol:clean]" + resetStyles,
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
			want:       "[style:clear]" + "[style:insertions][symbol:insertions]12",
		},
		{
			name:      "deletions",
			deletions: 12,
			want:      "[style:clear]" + "[style:deletions][symbol:deletions]12",
		},
		{
			name:       "insertions and deletions",
			insertions: 1,
			deletions:  2,
			want:       "[style:clear]" + "[style:insertions][symbol:insertions]1" + " " + "[style:deletions][symbol:deletions]2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{
					Styles: styles{
						Clear:      "[style:clear]",
						Deletions:  "[style:deletions]",
						Insertions: "[style:insertions]",
					},
					Symbols: symbols{
						Deletions:  "[symbol:deletions]",
						Insertions: "[symbol:insertions]",
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

func TestFlagsWithEmptySymbols(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "empty stashed symbol shows stash count (flags_without_count=false default)",
			styles: styles{
				Clear:    "[style:clear]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
			},
			symbols: symbols{
				Modified: "[symbol:mod]",
				Stashed:  "", // empty symbol should show count with default flags_without_count=false
			},
			st: &gitstatus.Status{
				NumStashed: 5,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
				},
			},
			want: "[style:clear]" + "[style:mod][symbol:mod]2 [style:stashed]5",
		},
		{
			name: "empty modified symbol shows modified count (flags_without_count=false default)",
			styles: styles{
				Clear:    "[style:clear]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
			},
			symbols: symbols{
				Modified: "", // empty symbol should show count with default flags_without_count=false
				Stashed:  "[symbol:stashed]",
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
				},
			},
			want: "[style:clear]" + "[style:mod]2 [style:stashed][symbol:stashed]1",
		},
		{
			name: "empty staged symbol shows staged count (flags_without_count=false default)",
			styles: styles{
				Clear:   "[style:clear]",
				Staged:  "[style:staged]",
				Stashed: "[style:stashed]",
			},
			symbols: symbols{
				Staged:  "", // empty symbol should show count with default flags_without_count=false
				Stashed: "[symbol:stashed]",
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumStaged: 3,
				},
			},
			want: "[style:clear]" + "[style:staged]3 [style:stashed][symbol:stashed]1",
		},
		{
			name: "empty untracked symbol shows untracked count (flags_without_count=false default)",
			styles: styles{
				Clear:     "[style:clear]",
				Untracked: "[style:untracked]",
				Stashed:   "[style:stashed]",
			},
			symbols: symbols{
				Untracked: "", // empty symbol should show count with default flags_without_count=false
				Stashed:   "[symbol:stashed]",
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumUntracked: 7,
				},
			},
			want: "[style:clear]" + "[style:stashed][symbol:stashed]1 [style:untracked]7",
		},
		{
			name: "empty conflict symbol shows conflict count (flags_without_count=false default)",
			styles: styles{
				Clear:    "[style:clear]",
				Conflict: "[style:conflict]",
				Stashed:  "[style:stashed]",
			},
			symbols: symbols{
				Conflict: "", // empty symbol should show count with default flags_without_count=false
				Stashed:  "[symbol:stashed]",
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumConflicts: 3,
				},
			},
			want: "[style:clear]" + "[style:conflict]3 [style:stashed][symbol:stashed]1",
		},
		{
			name: "empty clean symbol hides clean flag",
			styles: styles{
				Clear:   "[style:clear]",
				Clean:   "[style:clean]",
				Stashed: "[style:stashed]",
			},
			symbols: symbols{
				Clean:   "", // empty symbol should hide this flag
				Stashed: "[symbol:stashed]",
			},
			st: &gitstatus.Status{
				IsClean:    true,
				NumStashed: 1,
			},
			want: "[style:clear]" + "[style:stashed][symbol:stashed]1",
		},
		{
			name: "empty stashed symbol in clean state shows stash count (flags_without_count=false default)",
			styles: styles{
				Clear:   "[style:clear]",
				Clean:   "[style:clean]",
				Stashed: "[style:stashed]",
			},
			symbols: symbols{
				Clean:   "[symbol:clean]",
				Stashed: "", // empty symbol should show count with default flags_without_count=false
			},
			st: &gitstatus.Status{
				IsClean:    true,
				NumStashed: 1,
			},
			want: "[style:clear]" + "[style:stashed]1 [style:clean][symbol:clean]",
		},
		{
			name: "all symbols empty shows counts (flags_without_count=false default)",
			styles: styles{
				Clear:     "[style:clear]",
				Clean:     "[style:clean]",
				Staged:    "[style:staged]",
				Modified:  "[style:mod]",
				Conflict:  "[style:conflict]",
				Untracked: "[style:untracked]",
				Stashed:   "[style:stashed]",
			},
			symbols: symbols{
				Clean:     "",
				Staged:    "",
				Modified:  "",
				Conflict:  "",
				Untracked: "",
				Stashed:   "",
			},
			st: &gitstatus.Status{
				IsClean:    false,
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumStaged:    3,
					NumModified:  2,
					NumConflicts: 1,
					NumUntracked: 4,
				},
			},
			want: "[style:clear]" + "[style:staged]3 [style:conflict]1 [style:mod]2 [style:stashed]1 [style:untracked]4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols},
				st:     tt.st,
			}

			compareStrings(t, tt.want, f.flags())
		})
	}
}

func TestFlagsWithEmptySymbolsAndFlagsWithoutCount(t *testing.T) {
	tests := []struct {
		name    string
		styles  styles
		symbols symbols
		options options
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "empty symbols hide when flags_without_count=true",
			styles: styles{
				Clear:    "[style:clear]",
				Modified: "[style:mod]",
				Stashed:  "[style:stashed]",
			},
			symbols: symbols{
				Modified: "", // empty symbol should hide with flags_without_count=true
				Stashed:  "[symbol:stashed]",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
				},
			},
			want: "[style:clear]" + "[style:stashed][symbol:stashed]",
		},
		{
			name: "all empty symbols hide when flags_without_count=true",
			styles: styles{
				Clear:     "[style:clear]",
				Staged:    "[style:staged]",
				Modified:  "[style:mod]",
				Conflict:  "[style:conflict]",
				Untracked: "[style:untracked]",
				Stashed:   "[style:stashed]",
			},
			symbols: symbols{
				Staged:    "",
				Modified:  "",
				Conflict:  "",
				Untracked: "",
				Stashed:   "",
			},
			options: options{
				FlagsWithoutCount: true,
			},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumStaged:    3,
					NumModified:  2,
					NumConflicts: 1,
					NumUntracked: 4,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols, Options: tt.options},
				st:     tt.st,
			}

			compareStrings(t, tt.want, f.flags())
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

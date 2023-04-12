package tmux

import (
	"fmt"
	"io"
	"strings"
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
				Clear: "#[fg=default]",
				Clean: "#[fg=green,bold]",
			},
			symbols: symbols{
				Clean: "✔",
			},
			layout: []string{"branch", "..", "remote", "- ", "flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			want: "#[fg=default]" + "#[fg=green,bold]✔",
		},
		{
			name: "stash + clean flag",
			styles: styles{
				Clear:   "#[fg=default]",
				Clean:   "#[fg=green,bold]",
				Stashed: "#[fg=cyan,bold]",
			},
			symbols: symbols{
				Clean:   "✔",
				Stashed: "⚑ ",
			},
			layout: []string{"branch", "..", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				IsClean:    true,
				NumStashed: 1,
			},
			want: "#[fg=default]#[fg=cyan,bold]⚑ 1 #[fg=green,bold]✔",
		},
		{
			name: "mixed flags",
			styles: styles{
				Clear:    "#[fg=default]",
				Modified: "#[fg=red,bold]",
				Stashed:  "#[fg=cyan,bold]",
				Staged:   "[fg=green,bold]",
			},
			symbols: symbols{
				Modified: "✚ ",
				Stashed:  "⚑ ",
				Staged:   "● ",
			},
			layout: []string{"branch", "..", "remote", "- ", "flags"},
			st: &gitstatus.Status{
				NumStashed: 1,
				Porcelain: gitstatus.Porcelain{
					NumModified: 2,
					NumStaged:   3,
				},
			},
			want: "#[fg=default]" + "[fg=green,bold]● 3 #[fg=red,bold]✚ 2 #[fg=cyan,bold]⚑ 1",
		},
		{
			name: "mixed flags 2",
			styles: styles{
				Clear:     "#[fg=default]",
				Conflict:  "#[fg=red,bold]",
				Untracked: "#[fg=magenta,bold]",
			},
			symbols: symbols{
				Conflict:  "✖ ",
				Untracked: "… ",
			},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					NumConflicts: 42,
					NumUntracked: 17,
				},
			},
			want: "#[fg=default]" + "#[fg=red,bold]✖ 42 #[fg=magenta,bold]… 17",
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
		st      *gitstatus.Status
		want    string
	}{
		{
			name: "no divergence",
			styles: styles{
				Clear: "#[fg=default]",
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
				Clear: "#[fg=default]",
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
			want: "#[fg=default]" + "↓·4",
		},
		{
			name: "behind only",
			styles: styles{
				Clear: "#[fg=default]",
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
			want: "#[fg=default]" + "↑·12",
		},
		{
			name: "diverged both ways",
			styles: styles{
				Clear: "#[fg=default]",
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
			want: "#[fg=default]" + "↑·128↓·41",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{Styles: tt.styles, Symbols: tt.symbols},
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
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			compareStrings(t, tt.want, truncate(tt.s, tt.ellipsis, tt.max, tt.dir))
		})
	}
}

func TestFormatBranchAndRemote(t *testing.T) {
	st := gitstatus.Status{
		Porcelain: gitstatus.Porcelain{LocalBranch: "main"},
	}

	tests := []struct {
		layout []string
	}{
		{
			layout: []string{"branch", "...", "remote-branch"},
		},
	}

	for _, tt := range tests {
		f := &Formater{
			Config: Config{ /*Styles: tt.styles, Symbols: tt.symbols,*/ Layout: tt.layout /*Options: tt.options*/},
		}

		sb := &strings.Builder{}
		if err := f.Format(sb, &st); err != nil {
			t.Fatalf("Format error: %s", err)
			return
		}

		fmt.Println(sb.String())

		// compareStrings(t, tt.want, f.format())
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
				Clear:    "#[fg=default]",
				Clean:    "#[fg=green,bold]",
				Branch:   "#[fg=white,bold]",
				Modified: "#[fg=red,bold]",
				Remote:   "#[fg=cyan]",
			},
			symbols: symbols{
				Branch:   "⎇ ",
				Clean:    "✔",
				Modified: "✚ ",
			},
			layout: []string{"branch", " .. ", "remote", " - ", "flags"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "local",
					RemoteBranch: "remote",
					NumModified:  2,
				},
			},
			want: "#[fg=default]" + "#[fg=white,bold]⎇ " +
				"#[fg=default]" + "#[fg=white,bold]" + "local" +
				"#[fg=default]" + " .. " +
				"#[fg=default]" + "#[fg=cyan]remote" +
				"#[fg=default]" + " - " +
				"#[fg=default]" + "#[fg=red,bold]✚ 2" +
				resetStyles,
		},
		{
			name: "branch, different delimiter, flags",
			styles: styles{
				Clear:    "#[fg=default]",
				Branch:   "#[fg=white,bold]",
				Remote:   "#[fg=cyan]",
				Modified: "#[fg=red,bold]",
			},
			symbols: symbols{
				Branch:   "⎇ ",
				Ahead:    "↑·",
				Modified: "✚ ",
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
			want: "#[fg=default]" + "#[fg=white,bold]⎇ " +
				"#[fg=default]" + "#[fg=white,bold]" + "Local" +
				"#[fg=default]" + "~~" +
				"#[fg=default]" + "#[fg=red,bold]✚ 2" +
				resetStyles,
		},
		{
			name: "remote only",
			styles: styles{
				Clear:  "#[fg=default]",
				Branch: "#[fg=white,bold]",
				Remote: "#[fg=cyan]",
			},
			symbols: symbols{
				Branch: "⎇ ",
				Ahead:  "↑·",
			},
			layout: []string{"remote"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch:  "Local",
					RemoteBranch: "Remote",
					AheadCount:   1,
				},
			},
			want: "#[fg=default]" + "#[fg=cyan]Remote " +
				"#[fg=default]" + "↑·1" +
				resetStyles,
		},
		{
			name: "empty",
			styles: styles{
				Clear:    "#[fg=default]",
				Branch:   "#[fg=white,bold]",
				Modified: "#[fg=red,bold]",
			},
			symbols: symbols{
				Branch:   "⎇ ",
				Modified: "✚ ",
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
				Clear:  "#[fg=default]",
				Branch: "#[fg=white,bold]",
				Remote: "#[fg=cyan]",
			},
			symbols: symbols{
				Branch: "⎇ ",
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
			want: "#[fg=default]" + "#[fg=white,bold]" + "⎇ " +
				"#[fg=default]" + "#[fg=white,bold]" + "branchNa…" +
				"#[fg=default]" + "/" +
				"#[fg=default]" + "#[fg=cyan]" + "remote/b…" +
				resetStyles,
		},
		{
			name: "branch and remote, branch_max_len not zero and trim left",
			styles: styles{
				Clear:  "#[fg=default]",
				Branch: "#[fg=white,bold]",
				Remote: "#[fg=cyan]",
			},
			symbols: symbols{
				Branch: "⎇ ",
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
			want: "#[fg=default]" + "#[fg=white,bold]" + "⎇ " +
				"#[fg=default]" + "#[fg=white,bold]" + "...Branch " +
				"#[fg=default]" + "#[fg=cyan]" + "...Branch" +
				resetStyles,
		},
		{
			name: "issue-32",
			styles: styles{
				Clear:  "#[fg=default]",
				Branch: "#[fg=white,bold]",
			},
			symbols: symbols{
				Branch: "⎇ ",
			},
			layout: []string{"branch"},
			st: &gitstatus.Status{
				Porcelain: gitstatus.Porcelain{
					LocalBranch: "branchName",
				},
			},
			want: "#[fg=default]" + "#[fg=white,bold]" + "⎇ " +
				"#[fg=default]" + "#[fg=white,bold]" + "branchName" +
				resetStyles,
		},
		{
			name: "hide clean option true",
			styles: styles{
				Clear: "#[fg=default]",
				Clean: "#[fg=green,bold]",
			},
			symbols: symbols{
				Clean: "✔",
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
				Clear: "#[fg=default]",
				Clean: "#[fg=green,bold]",
			},
			symbols: symbols{
				Clean: "✔",
			},
			layout: []string{"flags"},
			st: &gitstatus.Status{
				IsClean: true,
			},
			options: options{
				HideClean: false,
			},
			want: "#[fg=default]" + "#[fg=green,bold]✔" + resetStyles,
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
			want:       "#[fg=default]" + "#[fg=green]Σ12",
		},
		{
			name:      "deletions",
			deletions: 12,
			want:      "#[fg=default]" + "#[fg=red]Δ12",
		},
		{
			name:       "insertions and deletions",
			insertions: 1,
			deletions:  2,
			want:       "#[fg=default]" + "#[fg=green]Σ1" + " " + "#[fg=red]Δ2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formater{
				Config: Config{
					Styles: styles{
						Clear:      "#[fg=default]",
						Deletions:  "#[fg=red]",
						Insertions: "#[fg=green]",
					},
					Symbols: symbols{
						Deletions:  "Δ",
						Insertions: "Σ",
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

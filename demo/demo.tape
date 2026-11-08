# gitmux demo, rendered with VHS: https://github.com/charmbracelet/vhs
#
# From the repository root:
#
#   vhs demo/demo.tape
#
# See demo/README.md for details.

Output demo/demo.gif

Require git
Require go
Require tmux

Set Shell "bash"
Set FontSize 22
Set Width 1200
Set Height 520
Set Padding 20
Set Theme "Catppuccin Frappe"
Set WindowBar Colorful
# Slow enough for viewers to follow what's being typed.
Set TypingSpeed 85ms

# Build gitmux, prepare the demo repository and start tmux, off camera.
Hide
Type "./demo/setup.sh && tmux -L gitmux-demo -f /tmp/gitmux-demo/tmux.conf new -s demo -c /tmp/gitmux-demo/project" Enter
Sleep 3s
Ctrl+L
Sleep 1s
Show

# Clean working tree, local branch tracking origin/main.
Sleep 3s

# Modified file.
Type "echo 'More docs.' >> README.md"
Sleep 500ms
Enter
Sleep 3s

# Untracked file.
Type "touch notes.txt"
Sleep 500ms
Enter
Sleep 3s

# Staged file.
Type "git add README.md"
Sleep 500ms
Enter
Sleep 3s

# Commit: ahead of origin/main.
Type "git commit -qm 'Update docs'"
Sleep 500ms
Enter
Sleep 3s

# Stash: the stash count shows, the working tree is clean.
Type "git stash -u -q"
Sleep 500ms
Enter
Sleep 3s

# Someone pushed to origin: local and remote branches diverged.
Type "git fetch -q"
Sleep 500ms
Enter
Sleep 3.5s

# Merge with conflicts: special state and conflict count.
Type "git merge feature"
Sleep 500ms
Enter
Sleep 4s

Type "git merge --abort"
Sleep 500ms
Enter
Sleep 3s

# New local branch, without upstream.
Type "git switch -q -c topic"
Sleep 500ms
Enter
Sleep 3s

# Outside of a Git repository, gitmux shows nothing.
Type "cd .."
Sleep 500ms
Enter
Sleep 3.5s

Hide
Type "tmux -L gitmux-demo kill-server" Enter
Sleep 500ms

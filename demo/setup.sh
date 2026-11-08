#!/usr/bin/env bash
#
# Prepares the environment used by demo.tape (see README.md in this
# directory). Builds gitmux from the current sources and creates, in $DEMO_DIR:
#   - bin/gitmux: gitmux binary
#   - origin.git: a bare repository acting as remote
#   - project:    a clone of origin, the demo happens there
#   - tmux.conf:  tmux configuration showing gitmux in the status bar
#   - bashrc:     minimal shell configuration (short prompt)
set -euo pipefail

DEMO_DIR=${DEMO_DIR:-/tmp/gitmux-demo}
ROOT=$(cd "$(dirname "$0")/.." && pwd)

rm -rf "$DEMO_DIR"
mkdir -p "$DEMO_DIR/bin"
(cd "$ROOT" && go build -o "$DEMO_DIR/bin/gitmux" .)

export GIT_AUTHOR_NAME=gitmux GIT_AUTHOR_EMAIL=demo@gitmux
export GIT_COMMITTER_NAME=gitmux GIT_COMMITTER_EMAIL=demo@gitmux
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1

cd "$DEMO_DIR"
git init -q --bare -b main origin.git
git clone -q origin.git project 2>/dev/null

cd project
git config user.name gitmux
git config user.email demo@gitmux
git config advice.mergeConflict false
git config advice.detachedHead false
printf '# My project\n\nSome documentation.\n' > README.md
printf 'package main\n\nfunc main() {}\n' > main.go
git add . && git commit -q -m "Initial commit"
git push -q -u origin main

# A local branch whose changes conflict with the ones made during the demo.
git switch -q -c feature
echo "Documentation from the feature branch." >> README.md
git commit -q -am "Document feature"
git switch -q main

# Someone else pushes to origin, so that 'git fetch' makes us behind.
git clone -q ../origin.git ../other 2>/dev/null
(
	cd ../other
	echo "// TODO" >> main.go
	git commit -q -am "Add TODO"
	git push -q origin main
)
rm -rf ../other

cat > "$DEMO_DIR/bashrc" <<'RC'
PS1='\[\e[1;34m\]\W\[\e[0m\] $ '
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
RC

cat > "$DEMO_DIR/tmux.conf" <<CONF
set -g default-command "bash --rcfile $DEMO_DIR/bashrc"
set -g status-interval 1
set -g status-style "bg=#303446,fg=#c6d0f5"
set -g status-left "#[bold] gitmux demo  "
set -g status-left-length 20
set -g status-right-length 120
set -g status-right '#($DEMO_DIR/bin/gitmux "#{pane_current_path}") '
set -g window-status-format ""
set -g window-status-current-format ""
CONF

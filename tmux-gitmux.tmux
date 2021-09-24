#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BINARY="$CURRENT_DIR/gitmux"
CONFIG="$CURRENT_DIR/gitmux.conf"

if [ ! -f "$BINARY" ]; then
  tmux split-window "cd $CURRENT_DIR && go build -o gitmux && echo 'Press any key to continue...' && read -k1"
fi

if [ ! -f "$CONFIG" ]; then
  ( cd "$CURRENT_DIR" ; ./gitmux -printcfg > gitmux.conf )
fi

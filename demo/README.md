# gitmux demo

`demo.tape` is a [VHS](https://github.com/charmbracelet/vhs) script that
records the gitmux demo GIF shown in the main README.

It builds gitmux from the current sources, creates a throwaway Git repository
in `/tmp/gitmux-demo` (see `setup.sh`), starts an isolated tmux server with
gitmux in its status bar, and walks through the main Git states: modified,
untracked and staged files, ahead/behind divergence, stash, merge conflict,
new branch, and leaving the repository.

## Requirements

`vhs` (which itself needs `ttyd` and `ffmpeg`), `go`, `git` and `tmux`.

## Recording

From the repository root:

    vhs demo/demo.tape

This writes `demo/demo.gif`. The README image is hosted on the `readme-images`
branch, so the generated GIF should be copied there.

Typing speed and pauses are set in `demo.tape`, so that viewers can read each
command and see the status bar update (tmux refreshes it every second, see
`status-interval` in `setup.sh`).

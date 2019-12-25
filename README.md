# Gitmux [![Build Status](https://travis-ci.com/arl/gitmux.svg?branch=master)](https://travis-ci.com/arl/gitmux) [![Go Report Card](https://goreportcard.com/badge/github.com/arl/gitmux)](https://goreportcard.com/report/github.com/arl/gitmux)


## **Gitmux** shows **Git** status in your **Tmux** status bar.

![Gitmux in action](https://raw.githubusercontent.com/arl/gitmux/readme-images/demo-small.gif)


## Description

Gitmux is a tmux addon that shows a minimal but useful **Git status** info in your tmux status bar.  
If the directory you're in is not a Git repository, **Gitmux** gets _out of your way_.

Many solutions already exist to keep an eye on Git status:
 - you can type git status each time you need it...we're too lazy for that!
 - you can embed git status into your shell prompt... that's overwhelming! I like to keep a small and tidy prompt.

And generally there's always a lot of empty space left in tmux status bar.

**Gitmux** comes with sensible defaults but you can customize everything: colors, symbols, information to show.

**To sum things up**:
 - you use **Tmux**
 - you're tired to type `git status`
 - you like a clean prompt

**Gitmux** might be just for you!


## Installation

* **Install a binary release for your platform** (preferred and simplest way) 

[Download](https://github.com/arl/gitmux/releases/latest) the latest binary.  
Add it to your `$PATH`.

* **Build from source**

Download and install the Go compiler from [golang.org](https://golang.org/dl/).  
Go get the latest source code, the dependencies, build and install all from one command:

```bash
go get -u github.com/arl/gitmux
```

## Usage

Simply add this line to your  `.tmux.conf`:

```
# Show Git working tree status
set -g status-right '#(gitmux #{pane_current_path})'
```


## Customize status string

Nothing simpler! First save `gitmux` config in a file:

```
gitmux -printcfg > .gitmux.conf
```

`gitmux` configuration is divided in 2 sections:
 - symbols are unicode characters
 - styles are tmux format strings (`man tmux` for reference)

Modify it, then feed that config each time you run `gitmux`:

```
gitmux -cfg .gitmux.conf
```

## Troubleshooting

If something goes wrong, please [file an issue](https://github.com/arl/gitmux/issues/new)
and indicate your tmux and gitmux versions,  
what you did, what you saw and what you expected yo see.  
Also you can run `gitmux -dbg` for debugging output.

```
tmux -V
gitmux -V
gitmux -dbg
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](./LICENSE)

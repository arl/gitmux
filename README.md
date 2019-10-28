# GiTmux [![Build Status](https://travis-ci.org/arl/gitmux.svg?branch=master)](https://travis-ci.org/arl/gitmux) [![Go Report Card](https://goreportcard.com/badge/github.com/arl/gitmux)](https://goreportcard.com/report/github.com/arl/gitstatus)

![GiTmux in action](https://raw.githubusercontent.com/arl/gitmux/readme-images/demo-small.gif)

**GiTmux** shows **Git** status in your **Tmux** status bar.

## Description

If the working directory is managed by Git, **GiTmux** will show **Git status**
information in a **minimal** and useful manner, right in Tmux status bar.  
GiTmux gets _out of your way_ when it has nothing to say (out of a Git
working tree).

**GiTmux** comes with sensible defaults though you can customize everything: colors, symbols, which information to show.

**To sum things up**:
 - you use **Tmux**
 - you're tired to type `git status`, or you're just _lazy_, like me
 - you want to keep your prompt tidy

then **GiTmux** is made for you!

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

Simply add this line to your  `.tmux.conf`

```
# Show Git working tree status
set -g status-right '#(gitmux -q -fmt tmux #{pane_current_path})'
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](./LICENSE)

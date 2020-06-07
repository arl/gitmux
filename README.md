<p align="center">
<img width="50%" height="50%" src="https://github.com/arl/gitmux/raw/readme-images/logo-transparent.png" />
</p>
<p align="center">Gitmux shows git status in your tmux status bar</p>
<hr>

<p align="center">
<a href="https://travis-ci.com/arl/gitmux">
  <img alt="travis-ci" src="https://travis-ci.com/arl/gitmux.svg?branch=master" />
</a>
<a href="https://goreportcard.com/report/github.com/arl/gitmux">
  <img alt="goreport" src="https://goreportcard.com/badge/github.com/arl/gitmux" />
</a>
<a href="https://opensource.org/licenses/MIT">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg" />
</a>
</p>

![demo](https://raw.githubusercontent.com/arl/gitmux/readme-images/demo-small.gif)


 - **easy**. Install it once and forget about it
 - **minimal**. Only show what you need when you need it
 - **discrete**. Disappear when current directory is not managed by Git
 - **shell-independent**. Work with sh, bash, zsh, fish, etc.
 - **highly configurable**. Colors, symbols and layout can be customized
 - **automatic**. Information auto-updates with respect to the current working directory

## Prerequisites

Works with all [tmux](https://github.com/tmux/tmux) versions.

## Installing

### Binary release
[Download the latest](https://github.com/arl/gitmux/releases/latest) binary for your platform/architecture and uncompress it.

### From source

[Download and install a Go compiler](https://golang.org/dl/) (Go 1.10 or later).
Run `go get` to build and install `gitmux`:

    go get -u github.com/arl/gitmux

## Getting started

Add this line to your  `.tmux.conf`:

    set -g status-right '#(gitmux "#{pane_current_path}")'


## Customizing

`gitmux` output can be customized via a configuration file in YAML format.

The gitmux configuration file is in YAML format:

```yaml
tmux:
  symbols:
    branch: '⎇ '
    hashprefix: ':'
    ahead: ↑·
    behind: ↓·
    staged: '● '
    conflict: '✖ '
    modified: '✚ '
    untracked: '… '
    stashed: '⚑ '
    clean: ✔
  styles:
    state: '#[fg=red,bold]'
    branch: '#[fg=white,bold]'
    remote: '#[fg=cyan]'
    staged: '#[fg=green,bold]'
    conflict: '#[fg=red,bold]'
    modified: '#[fg=red,bold]'
    untracked: '#[fg=magenta,bold]'
    stashed: '#[fg=cyan,bold]'
    clean: '#[fg=green,bold]'
  layout: [branch, .., remote, ' - ', flags]
```

First, save the default configuration to a new file:

    gitmux -printcfg > .gitmux.conf

Modify the line in `.tmux.conf`, passing the path of the configuration file as argument to `gitmux`

    gitmux -cfg .gitmux.conf

Open `.gitmux.conf` and modify it, replacing symbols, styles and layout to suit your needs.

In `tmux` status bar, `gitmux` output immediately reflects the changes you make to the configuration.

`gitmux` configuration is split into 3 sections:
 - `symbols`: they're just strings of unicode characters
 - `styles`: tmux format strings
 - `layout`: list of `gitmux` layout components, defines the component to show and in their order.


### Symbols

```yaml
  symbols:
    branch: '⎇ '      # shown before `branch`
    hashprefix: ':'    # shown before a Git hash (in 'detached HEAD' state)
    ahead: ↑·          # shown before 'ahead count' when local/remote branch diverges`
    behind: ↓·         # shown before 'behind count' when local/remote branch diverges`
    staged: '● '       # shown before the 'staged files' count
    conflict: '✖ '     # shown before the 'conflicts' count
    modified: '✚ '     # shown before the 'modified files' count
    untracked: '… '    # shown before the 'untracked files' count
    stashed: '⚑ '      # shown before the 'stash' count
    clean: ✔           # shown when the working tree is clean
```


### Styles

Styles are tmux format strings. For full reference, search for `message-command-style` 
in `tmux` manual `man tmux`.


### Layout components

This is the list of the possible components of the `layout`:

| Layout Component |                 Description                 |         Example        |
|:----------------:|---------------------------------------------|:----------------------:|
| `branch`         | Local branch name                           |        `master`        |
| `remote`         | Remote branch name                          |     `origin/master`    |
| `divergence`     | Divergence local/remote branch, if any      |        `↓·2↑·1`        |
| `remote`         | Alias for `remote-branch divergence`        | `origin/master ↓·2↑·1` |
| `flags`          | Symbols representing the working tree state |      `✚ 1 ⚑ 1 … 2`     |
| any string `foo` | Any other string is directly shown          |          `foo`         |



Example layouts:
```
layout: [branch, '..', remote, ' - ', flags]
layout: [branch, '..', remote-flags, divergence, ' - ', flags]
layout: [branch]
layout: [flags, branch]
layout: [flags, ~~~, branch]
```

## Troubleshooting

Please report anything by [filing an issue](https://github.com/arl/gitmux/issues/new).


## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License: [MIT](./LICENSE)

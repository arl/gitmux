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


 - **easy**. Install and forget about it
 - **minimal**. It shows what you need when you need it
 - **discrete**. Disappears if the current directory is not part of a Git tree
 - **shell-independent**. Works with all shells bash, zsh, fish, whateversh
 - **customizable**. Colors, symbols and layout can be customized

## Prerequisites

Works with all decently recent [tmux](https://github.com/tmux/tmux) versions.

## Installing

### Binary release
[Download the latest](https://github.com/arl/gitmux/releases/latest) binary for your platform/architecture and uncompress it.

### From source

[Download and install a Go compiler](https://golang.org/dl/) (Go 1.10 or later).
Run `go get` to build and install `gitmux`:

    go get -u github.com/arl/gitmux

## Getting started

If your `tmux` version supports `pane_current_pane` (tmux v2.1+),
just add this line to your `.tmux.conf`:

    set -g status-right '#(gitmux "#{pane_current_path}")'

If your `tmux` doesn't support `pane_current_path` then you can use 
a [bash-specific solution](https://github.com/arl/gitmux/issues/19#issuecomment-594735939)
to achieve relatively similar behaviour: `gitmux` will refresh after every shell command 
you run or when you switch windows, however it won't refresh automatically, nor when switching panes.  

Note that `tmux v2.1` was released in 2015 so you're probably better off updating to a more recent version anyway 🙂.

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
    clear: '#[fg=default]'
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
  options:
    branch_max_len: 0
```

First, save the default configuration to a new file:

    gitmux -printcfg > .gitmux.conf

Modify the line in `.tmux.conf`, passing the path of the configuration file as argument to `gitmux`

    gitmux -cfg .gitmux.conf

Open `.gitmux.conf` and modify it, replacing symbols, styles and layout to suit your needs.

In `tmux` status bar, `gitmux` output immediately reflects the changes you make to the configuration.

`gitmux` configuration is split into 4 sections:
 - `symbols`: they're just strings of unicode characters
 - `styles`: tmux format strings
 - `layout`: list of `gitmux` layout components, defines the component to show and in their order.
 - `options`: additional configuration options


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

Styles are tmux format strings used to specify text colors and attributes.
See [`tmux` styles reference](https://man7.org/linux/man-pages/man1/tmux.1.html#STYLES).


### Layout components

This is the list of the possible components of the `layout`:

| Layout component |                 Description                               |        Example         |
|:----------------:|:----------------------------------------------------------|:----------------------:|
| `branch`         | local branch name                                         |        `master`        |
| `remote-branch`  | remote branch name                                        |     `origin/master`    |
| `divergence`     | divergence local/remote branch, if any                    |        `↓·2↑·1`        |
| `remote`         | alias for `remote-branch` followed by `divergence`        | `origin/master ↓·2↑·1` |
| `flags`          | Symbols representing the working tree state               |      `✚ 1 ⚑ 1 … 2`     |
| any string `foo` | Any other string is directly shown                        |          `foo`         |



Example layouts:
```
layout: [branch, '..', remote, ' - ', flags]
layout: [branch, '..', remote-flags, divergence, ' - ', flags]
layout: [branch]
layout: [flags, branch]
layout: [flags, ~~~, branch]
```


### Additional options

This is the list of additional configuration `options`:

| Option           | Description                                                | Default        |
|:-----------------|:-----------------------------------------------------------|:---------------|
| `branch_max_len` | Maximum displayed length for local and remote branch names | `0` (no limit) |


## Troubleshooting

Please report anything by [filing an issue](https://github.com/arl/gitmux/issues/new).


## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License: [MIT](./LICENSE)

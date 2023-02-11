<p align="center">
<img width="50%" height="50%" src="https://github.com/arl/gitmux/raw/readme-images/logo-transparent.png" />
</p>
<p align="center">Gitmux shows git status in your tmux status bar</p>
<hr>

<p align="center">

<a href="https://github.com/arl/gitmux/actions/workflows/ci.yml">
  <img alt="tests" src="https://github.com/arl/gitmux/actions/workflows/ci.yml/badge.svg" />
</a>
<a href="https://github.com/arl/gitmux/actions/workflows/cd.yml">
  <img alt="tests" src="https://github.com/arl/gitmux/actions/workflows/cd.yml/badge.svg" />
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
 - **minimal**. Just shows what you need, when you need it
 - **discrete**. Get out of your way if current directory is not in a Git tree
 - **shell-agnostic**. Does not rely on shell-features so works with all of them
 - **customizable**. Colors, symbols and layout are configurable


## Prerequisites

Works with all decently recent [tmux](https://github.com/tmux/tmux) versions.

## Installing

### Binary release

[Download the latest](https://github.com/arl/gitmux/releases/latest) binary for your platform/architecture and uncompress it.


### Homebrew (tap) macOS and linux, amd64 and arm64

Install the latest version with:

```sh
brew tap arl/arl
brew install gitmux
```

### AUR

Arch Linux users can download the [gitmux](https://aur.archlinux.org/packages/gitmux), [gitmux-bin](https://aur.archlinux.org/packages/gitmux-bin) or [gitmux-git](https://aur.archlinux.org/packages/gitmux-git) AUR package.

### From source

[Download and install a Go compiler](https://golang.org/dl/) (Go 1.16 or later).
Run `go install` to build and install `gitmux`:

```bash
go install github.com/arl/gitmux@latest
```

## Getting started

If your `tmux` version supports `pane_current_path` (tmux v2.1+),
just add this line to your `.tmux.conf`:

    set -g status-right '#(gitmux "#{pane_current_path}")'

If your `tmux` doesn't support `pane_current_path` then you can use 
a [bash-specific solution](https://github.com/arl/gitmux/issues/19#issuecomment-594735939)
to achieve relatively similar behaviour: `gitmux` will refresh after every shell command 
you run or when you switch windows, however it won't refresh automatically, nor when switching panes.  

Note that `tmux v2.1` was released in 2015 so you're probably better off updating to a more recent version anyway 🙂.

## Customizing

`gitmux` output can be customized via a configuration file in YAML format.

This is the default gitmux configuration file, in YAML format:

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
        divergence: '#[fg=default]'
    layout: [branch, .., remote-branch, divergence, ' - ', flags]
    options:
        branch_max_len: 0
        branch_trim: right
```

First, save the default configuration to a new file:

    gitmux -printcfg > .gitmux.conf

Modify the line in `.tmux.conf`, passing the path of the configuration file as argument to `gitmux` via the `-cfg` flag

    set -g status-right '#(gitmux -cfg .gitmux.conf "#{pane_current_path}")'

Open `.gitmux.conf` and modify it, replacing symbols, styles and layout to suit your needs.

In `tmux` status bar, `gitmux` output immediately reflects the changes you make to the configuration.

`gitmux` configuration is split into 4 sections:
 - `symbols`: they're just strings of unicode characters
 - `styles`: tmux format strings
 - `layout`: list of `gitmux` layout components, defines the component to show and in their order.
 - `options`: additional configuration options


### Symbols

The `symbols` section describes the symbols `gitmux` prints for the various components of the status string.

```yaml
  symbols:
    branch: '⎇ '      # Shown before a branch
    hashprefix: ':'    # Shown before a Git hash (in 'detached HEAD' state)
    ahead: ↑·          # Shown before the 'ahead count' when local and remote branch diverge
    behind: ↓·         # Shown before the 'behind count' when local/remote branch diverge
    staged: '● '       # Shown before the 'staged' files count
    conflict: '✖ '     # Shown before the 'conflicts' count
    modified: '✚ '     # Shown before the 'modified' files count
    untracked: '… '    # Shown before the 'untracked' files count
    stashed: '⚑ '      # Shown before the 'stash' count
    clean: ✔           # Shown when the working tree is clean (empty staging area)
```


### Styles

Styles are tmux format strings used to specify text colors and attributes.
See [`tmux` styles reference](https://man7.org/linux/man-pages/man1/tmux.1.html#STYLES).

```yaml
  styles:
    clear: '#[fg=default]'          # Style clearing previous styles (printed before each component)
    state: '#[fg=red,bold]'         # Style of the special states strings like [rebase], [merge], etc.
    branch: '#[fg=white,bold]'      # Style of the local branch name
    remote: '#[fg=cyan]'            # Style of the remote branch name
    divergence: "#[fg=yellow]"      # Style of the 'divergence' string
    staged: '#[fg=green,bold]'      # Style of the 'staged' files count
    conflict: '#[fg=red,bold]'      # Style of the 'conflicts' count
    modified: '#[fg=red,bold]'      # Style of the 'modified' files count
    untracked: '#[fg=magenta,bold]' # Style of the 'modified' files count
    stashed: '#[fg=cyan,bold]'      # Style of the 'stash' entries count
    clean: '#[fg=green,bold]'       # Style of the 'clean' symbol
```

### Layout components

The layout is a list of the components `gitmux` shows, and the order in
which they appear on `tmux` status bar.

For example, the default `gitmux` layout shows is:

```yaml
layout: [branch, .., remote-branch, divergence, " - ", flags]
```

It shows, in that order:
 - the local branch name,
 - 2 dots characters `..`,
 - the remote branch name
 - the local/remote divergence
 - a `-` character
 - and finally the flags representing the working tree state

Note that elements only appear when they make sense, for example if local and
remote branch are aligned, the divergence string won't show up. Same thing for
the remote branch, etc.

But you can anyway choose to never show some components if you wish, or to present
them in a different order.

This is the list of the possible components of the `layout`:

| Layout component | Description                                        |       Example        |
| :--------------: | :------------------------------------------------- | :------------------: |
|     `branch`     | local branch name                                  |        `main`        |
| `remote-branch`  | remote branch name                                 |    `origin/main`     |
|   `divergence`   | divergence local/remote branch, if any             |       `↓·2↑·1`       |
|     `remote`     | alias for `remote-branch` followed by `divergence` | `origin/main ↓·2↑·1` |
|     `flags`      | Symbols representing the working tree state        |    `✚ 1 ⚑ 1 … 2`     |
| any string `foo` | Any other string is directly shown                 |        `foo`         |


Some example layouts:

 - default layout:

```yaml
layout: [branch, .., remote-branch, divergence, " - ", flags]
```

 - some more minimal layouts:

```yaml
layout: [branch, divergence, " - ", flags]
```
```yaml
layout: [flags, " ", branch]
```


### Additional options

This is the list of additional configuration `options`:

| Option           | Description                                                | Default            |
| :--------------- | :--------------------------------------------------------- | :----------------- |
| `branch_max_len` | Maximum displayed length for local and remote branch names | `0` (no limit)     |
| `branch_trim`    | Trim left or right end of the branch (`right` or `left`)   | `right` (trailing) |


## Troubleshooting

Check the opened and closed issues and don't hesitate to report anything by [filing a new one](https://github.com/arl/gitmux/issues/new). 


### Gitmux takes too long to refresh?

In case gitmux takes too long to refresh, try to decrease the value of the `status-interval` option.
A reasonable value is 2 seconds, which you can set in `.tmux.conf` with:

    set -g status-interval 2

Check out [tmux man page](https://www.man7.org/linux/man-pages/man1/tmux.1.html#OPTIONS) for more details.


## Contributing

Pull requests are welcome.  
For major changes, please open an issue first to open a discussion.


## License: [MIT](./LICENSE)

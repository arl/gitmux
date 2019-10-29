# Gitmux

Gitmux shows Git in your tmux status bar.

![Gitmux in action](https://raw.githubusercontent.com/arl/gitmux/readme-images/demo-small.gif)

## Installation

### Install a binary release for your platform

[Download](https://github.com/arl/gitmux/releases/latest) the latest pre-compiled binary release.  
Add it to your `$PATH`.

### Or, build from source

[Download](https://golang.org/dl/) the latest Go binary.  
Get **gitmux** source and and build it:

```bash
go get github.com/arl/gitmux
```

## Usage

Add this line to your  `.tmux.conf`

```
# Show Git working tree status
set -g status-right '#(gitmux -q -fmt tmux #{pane_current_path})'
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](./LICENSE)

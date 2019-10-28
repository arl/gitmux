# Gitmux

Gitmux shows Git in your tmux status bar.

## Installation

### Pre-compiled binaries for all supported platforms

Download the latest binary release for your platform and add it to your `$PATH`.

### Build from source

Install the latest [Go version](https://golang.org/dl/) and then run:

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
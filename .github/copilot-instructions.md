This is a Go based repository buliding a command line tool that converts the status of a git working tree for a specific directory into a tmux status string the end user can add to their tmux status line. It is designed to be used with the `tmux` terminal multiplexer.

## Code Standards

### Required Before Each Commit

- Run `go test ./...` before committing any changes to ensure tests pass
- Go code should be formatted with `gofmt`

## Repository Structure

- `./`: repository root, main package, default configuration file, and main entry point
- `tmux/`: tmux formatting
- `json/`: used by gitmux -dbg to print the git working tree status as a json object, for debugging purposes
- `testdata/`: testscripts fixtures when actual gitmux output is checked against some specific conditions.

## Key Guidelines

1. Follow Go best practices and idiomatic patterns
2. Maintain existing code structure and organization
3. Write unit tests for new functionality. Use table-driven unit tests when possible.
4. Document public options
5. Keep table aligned in the README.md file, for a nice experience even when reading the file in a terminal

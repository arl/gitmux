#!/bin/sh

_have() { type "$1" >/dev/null 2>&1; }

help() {
    cat <<EOF
Usage: ${0##*/} [options] [version]

Options:
  -h, --help      Show this help message and exit.

Arguments:
  version         Specify the version of gitmux to install. If not provided, the latest version will be installed.

Environment Variables:
  INSTALL_PATH    Specify the directory where gitmux should be installed. Required

Description:
  This script installs the gitmux tool, which provides Git status information in the tmux status line.
  It automatically detects the system architecture and downloads the appropriate version of gitmux.

Examples:
  INSTALL_PATH=/usr/local/bin ${0##*/}            Install the latest version of gitmux to /usr/local/bin.
  INSTALL_PATH=/usr/local/bin ${0##*/} v0.7.0     Install version 0.7.0 of gitmux to /usr/local/bin.

EOF
}

gh_install() {
    ver="$1"
    repo='arl/gitmux'
    [ -z "$ver" ] && {
        latest="https://api.github.com/repos/$repo/releases/latest"
        ver=$(curl -sS "$latest" | grep tarball_url | sed 's>.*: "\(.*\)".*>\1>') && test -n "$ver"
        ver=${ver##*/}
    }
    tarname=''
    archi=$(uname -sm)
    case "$archi" in
    Darwin\ arm64) tarname="gitmux_${ver}_macOS_arm64.tar.gz" ;;
    Darwin\ x86_64) tarname="gitmux_${ver}_macOS_amd64.tar.gz" ;;
    Linux\ aarch64*) tarname="gitmux_${ver}_linux_arm64.tar.gz" ;;
    Linux\ *64) tarname="gitmux_${ver}_linux_amd64.tar.gz" ;;
    *) echo "Unsupported architecture" && return 1 ;;
    esac
    tmpdir="$(mktemp -d)"
    cd "$tmpdir" || :
    curl -sSLO "https://github.com/$repo/releases/download/$ver/$tarname"
    tar -xf gitmux*.tar.gz &&
        mv gitmux "${INSTALL_PATH}"
    [ -d "$tmpdir" ] && rm -rf "$tmpdir"
}

if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then help && exit 0; fi

! _have curl && echo "This script depends on curl" && exit 1
[ -z "$INSTALL_PATH" ] && echo "Please set the INSTALL_PATH envvar to specify installation directory" && exit 1

echo "Installing gitmux..."
if gh_install "$@"; then
    echo "Successfully installed gitmux"
else
    echo "Could not install gitmux" && exit 1
fi

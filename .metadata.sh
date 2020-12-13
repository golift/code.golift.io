# This file is only used for automatic builds.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
# Powered by Application Builder: https://github.com/golift/application-builder

# Bring in dynamic repo/pull/source info.
source $(dirname "${BASH_SOURCE[0]}")/init/buildinfo.sh

# Must match the repo name to make things easy. Otherwise, fix some other paths.
BINARY="turbovanityurls"
REPO="turbovanityurls"
# github username
GHUSER="golift"
# Github repo containing homebrew formula repo.
HBREPO="golift/homebrew-mugs"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift"
DESC="HTTP Server providing vanity go import paths."
# FIX THIS.
GOLANGCI_LINT_ARGS="--enable-all -D gochecknoglobals -D lll -D unused -D wsl \
  -D wrapcheck -D paralleltest -D testpackage -D nlreturn -D gosimple \
  -D exhaustivestruct -D noctx -D gosec -D goerr113 -D funlen"
# Example must exist at examples/$CONFIG_FILE.example
CONFIG_FILE="config.yaml"
LICENSE="Apache-2.0"
# FORMULA is either 'service' or 'tool'. Services run as a daemon, tools do not.
# This affects the homebrew formula (launchd) and linux packages (systemd).
FORMULA="service"

# Used for source links and wiki links.
SOURCE_URL="https://github.com/${GHUSER}/${REPO}"
# Used for documentation links.
URL="${SOURCE_URL}"

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
# This makes a path like github.com/user/hello-world/helloworld.Version=1.3.3
# Name the Version-containing library the same as the github repo, without dashes.
VERSION_PATH="main"

# Used by homebrew downloads.
SOURCE_PATH=https://codeload.github.com/${GHUSER}/${REPO}/tar.gz/v${VERSION}

export BINARY GHUSER HBREPO MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE
export LICENSE FORMULA SOURCE_URL URL VERSION_PATH SOURCE_PATH

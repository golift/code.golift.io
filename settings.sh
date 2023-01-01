# Used by Docker, homebrew and aur (arch) builds to inject version and build info.
# This same logic exists in the Makefile too (for linux packages).
# This file must be simple and work with BASH and MAKE.

MAINT="David Newhall II <captain at golift dot io>"
DESC="HTTP Server providing vanity go import paths."
LICENSE="Apache-2.0"
SOURCE_URL="https://github.com/golift/turbovanityurls"
VENDOR="Go Lift <code@golift.io>"

DATE="$(date -u +%Y-%m-%dT%H:%M:00Z)"
VERSION=$(git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1) | tr -d v)
[ "$VERSION" != "" ] || VERSION=development
# This produces a 0 in some environments (like Homebrew), but it's only used for packages.
ITERATION=$(git rev-list --count --all || echo 0)
COMMIT="$(git rev-parse --short HEAD || echo 0)"
BRANCH="$(git rev-parse --abbrev-ref HEAD || echo ${GITHUB_REF_NAME})"

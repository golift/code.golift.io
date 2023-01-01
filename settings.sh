# This file is only used for automatic builds.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
# Powered by Application Builder: https://github.com/golift/application-builder

# Bring in dynamic repo/pull/source info.
MAINT="David Newhall II <captain at golift dot io>"
DESC="HTTP Server providing vanity go import paths."
LICENSE="Apache-2.0"
SOURCE_URL="https://github.com/golift/turbovanityurls"
VENDOR="Go Lift <code@golift.io>"
export MAINT DESC LICENSE SOURCE_URL VENDOR

DATE="$(date -u +%Y-%m-%dT%H:%M:00Z)"
VERSION=$(git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1) | tr -d v)
[ "$VERSION" != "" ] || VERSION=development
# This produces a 0 in some environments (like Homebrew), but it's only used for packages.
ITERATION=$(git rev-list --count --all || echo 0)
COMMIT="$(git rev-parse --short HEAD || echo 0)"
GIT_BRANCH="$(git rev-parse --abbrev-ref HEAD || echo unknown)"
BRANCH="${GIT_BRANCH:-${GITHUB_REF_NAME}}"
export DATE VERSION ITERATION COMMIT BRANCH

export GHUSER MAINT DESC CONFIG_FILE
export LICENSE FORMULA SOURCE_URL SOURCE_PATH

### Optional ###

# Import this signing key only if it's in the keyring.
gpg --list-keys 2>/dev/null | grep -q B93DD66EF98E54E2EAE025BA0166AD34ABC5A57C
[ "$?" != "0" ] || export SIGNING_KEY=B93DD66EF98E54E2EAE025BA0166AD34ABC5A57C

export WINDOWS_LDFLAGS=""
export MACAPP=""
export EXTRA_FPM_FLAGS=""

# Make sure Docker builds work locally.
# These do not affect automated builds, just allow the docker build scripts to run from a local clone.
[ -n "$SOURCE_BRANCH" ] || export SOURCE_BRANCH=$BRANCH
[ -n "$DOCKER_TAG" ] || export DOCKER_TAG=$(echo $SOURCE_BRANCH | sed 's/^v*\([0-9].*\)/\1/')
[ -n "$DOCKER_REPO" ] || export DOCKER_REPO="golift/turbovanityurls"
[ -n "$IMAGE_NAME" ] || export IMAGE_NAME="${DOCKER_REPO}:${DOCKER_TAG}"
[ -n "$DOCKERFILE_PATH" ] || export DOCKERFILE_PATH="init/docker/Dockerfile"
#!/bin/bash -x

# Deploys a new aur PKGBUILD file to am arch linux aur github repo.
# Run by GitHub Actions when a new release is created on GitHub.

source settings.sh

SOURCE_PATH=https://github.com/golift/turbovanityurls/archive/v${VERSION}.tar.gz
echo "==> Using URL: $SOURCE_PATH"
SHA256=$(curl -sL $SOURCE_PATH | openssl dgst -r -sha256 | awk '{print $1}')

push_it() {
  pushd release_repo
  git add .
  git commit -m "Update turbovanityurls on Release: v${VERSION}-${ITERATION}"
  git push
  popd
  rm -rf release_repo
}

# Make an id_rsa file with our secret.
mkdir -p $HOME/.ssh
KEY_FILE="$(mktemp -u $HOME/.ssh/XXXXX)"
echo "${DEPLOY_KEY}" > "${KEY_FILE}"
chmod 600 "${KEY_FILE}"
# Configure ssh to use this secret on a custom github hostname.
GITHUB_HOST="github.$(basename $KEY_FILE)"
printf "%s\n" \
  "Host $GITHUB_HOST" \
  "  HostName github.com" \
  "  IdentityFile ${KEY_FILE}" \
  "  StrictHostKeyChecking no" \
  "  LogLevel ERROR" | tee -a $HOME/.ssh/config

git config --global user.email "turbovanityurls@auto.releaser"
git config --global user.name "turbovanityurls-auto-releaser"

rm -rf release_repo
git clone git@${GITHUB_HOST}:golift/aur.git release_repo
mkdir -p release_repo/turbovanityurls

sed -e "s/{{VERSION}}/${VERSION}/g" \
    -e "s/{{Iter}}/${ITERATION}/g" \
    -e "s/{{SHA256}}/${SHA256}/g" \
    -e "s/{{Desc}}/${DESC}/g" \
    -e "s%{{SOURCE_URL}}%${SOURCE_URL}%g" \
    -e "s%{{SOURCE_PATH}}%${SOURCE_PATH}%g" \
    init/archlinux/PKGBUILD.template | tee release_repo/turbovanityurls/PKGBUILD

sed -e "s/{{VERSION}}/${VERSION}/g" \
    -e "s/{{Iter}}/${ITERATION}/g" \
    -e "s/{{SHA256}}/${SHA256}/g" \
    -e "s/{{Desc}}/${DESC}/g" \
    -e "s%{{SOURCE_URL}}%${SOURCE_URL}%g" \
    -e "s%{{SOURCE_PATH}}%${SOURCE_PATH}%g" \
    init/archlinux/SRCINFO.template | tee release_repo/turbovanityurls/.SRCINFO

tee release_repo/turbovanityurls/turbovanityurls.aur.install << EOF
post_upgrade() {
  /bin/systemctl restart turbovanityurls
}

pre_remove() {
  /bin/systemctl stop turbovanityurls
  /bin/systemctl disable turbovanityurls
}
EOF

push_it

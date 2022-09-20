#!/bin/bash
#
# This script installs the Go-Lift APT and/or YUM repo(s) on a Linux system.
# Package Repository Hosting Provided by: https://PackageCloud.io
# When run on macOS it attempts to tap the golift homebrew repo.
# Optionally triggers a package install if $1 is non-empty.
#
### Install hello-world:
# curl -sL https://golift.io/repo.sh | sudo bash -s - hello-world
#
### Just install the repo:
# curl -sL https://golift.io/repo.sh | sudo bash
#

APT=$(which apt)
YUM=$(which yum)
BREW=$(which brew)
PKG=$1
extra=""

if [ -d /usr/share/keyrings ]; then 
  curl -sL https://packagecloud.io/golift/pkgs/gpgkey | gpg --dearmor > /usr/share/keyrings/golift-archive-keyring.gpg
  extra="[signed-by=/usr/share/keyrings/golift-archive-keyring.gpg] "
fi

# All Debian/Ubuntu/etc packages are in the ubuntu/focal repo.
###
if [ -d /etc/apt/sources.list.d ] && [ "$APT" != "" ]; then
  curl -sL https://packagecloud.io/golift/pkgs/gpgkey | apt-key add -
  echo "deb ${extra}https://packagecloud.io/golift/pkgs/ubuntu focal main" > /etc/apt/sources.list.d/golift.list
  apt update
  [ "$PKG" = "" ] || apt install $PKG
fi

# All RedHat/CentOS/etc packages are in the el/6 repo.
###
if [ -d /etc/yum.repos.d ] && [ "$YUM" != "" ]; then
  cat <<EOF > /etc/yum.repos.d/golift.repo
[golift]
name=golift
baseurl=https://packagecloud.io/golift/pkgs/el/6/\$basearch
repo_gpgcheck=1
gpgcheck=1
enabled=1
gpgkey=https://packagecloud.io/golift/pkgs/gpgkey
       https://packagecloud.io/golift/pkgs/gpgkey/golift-pkgs-7F7791485BF8996D.pub.gpg
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
metadata_expire=300
EOF

  yum -q makecache -y --disablerepo='*' --enablerepo='golift'
  [ "$PKG" = "" ] || yum install $PKG
fi

# All macOS packages are in the same homebrew repo.
###
if [ "$(uname -s 2>/dev/null)" = "Darwin" ] && [ "$BREW" != "" ]; then
  brew tap golift/mugs
  [ "$PKG" = "" ] || brew install $PKG
fi

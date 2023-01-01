//go:build !windows && !darwin && !freebsd

package service

// DefaultConfFile is where to find config if -c is not provided.
const DefaultConfFile = `/etc/turbovanityurls/config.yaml`

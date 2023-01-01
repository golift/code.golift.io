#!/bin/sh
#
# FreeBSD rc.d startup script for turbovanityurls.
#
# PROVIDE: turbovanityurls
# REQUIRE: networking syslog
# KEYWORD:

. /etc/rc.subr

name="turbovanityurls"
real_name="turbovanityurls"
rcvar="turbovanityurls_enable"
turbovanityurls_command="/usr/local/bin/${real_name}"
turbovanityurls_user="turbovanityurls"
turbovanityurls_config="/usr/local/etc/${real_name}/config.yaml"
pidfile="/var/run/${real_name}/pid"

# This runs `daemon` as the `turbovanityurls_user` user.
command="/usr/sbin/daemon"
command_args="-P ${pidfile} -r -t ${real_name} -T ${real_name} -l daemon ${turbovanityurls_command} -c ${turbovanityurls_config}"

load_rc_config ${name}
: ${turbovanityurls_enable:=no}

# Make a place for the pid file.
mkdir -p $(dirname ${pidfile})
chown -R $turbovanityurls_user $(dirname ${pidfile})

# Suck in optional exported override variables.
# ie. add something like the following to this file: export UP_POLLER_DEBUG=true
[ -f "/usr/local/etc/defaults/${real_name}" ] && . "/usr/local/etc/defaults/${real_name}"

# Go!
run_rc_command "$1"

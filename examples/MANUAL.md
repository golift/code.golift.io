turbovanityurls(1) -- Provide vanity go import paths
===

SYNOPSIS
---
`turbovanityurls -c /etc/turbovanityurls/config.yaml`

This daemon prints hello world.

OPTIONS
---
`turbovanityurls [-c <config-file>] [-h] [-v]`

    -c <config-file>
        Provide a configuration file (instead of the default).
        The default is ./config.yaml, but this may change in the future.

    -l <listen-addr>
        Provide a listen address to bind to. Default is :$PORT. PORT is taken
        from the environment. If PORT is unset the default is :8080.

    -v
        Display version and exit.

    -h
        Display usage and exit.

CONFIGURATION
---

`Config File Parameters`

..todo..

AUTHOR
---
*   GoogleCloudPlatform - 2017-2018
*   David Newhall II    - 2019

LOCATION
---
*   Turbo Vanity URLs: [https://github.com/golift/turbovanityurls](https://github.com/golift/turbovanityurls)

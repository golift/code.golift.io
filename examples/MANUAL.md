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

    title
      Used as the page and html title on the Index page.

    host
      Used as the import path host. Recommend setting this.

    cache_age                   default: 86400
      Cache-Control header max-age value. This is how long to tell upstream proxy
      servers they may cache our vanity pages for.

    src
      Used as a <(source)> link on the index page if set.

    links                       list
      Links are displayed in a list on the index page. Each link has a title and a url.

    bd_path
      This parameter is used to control badgedata. Badgedata is a custom library
      that provides "data" for badges. This is a feature used by golift.io, and
      most people will probably disable this. Set it to "" or remove the line from
      your config to disable badge data.

    redir_paths                 list
      These values are used in a string match to check it a path can be redirected.
      This only works if a path has `redir` set to a non-empty value. If the request
      URI contains one of these values it will be redirected. This setting is global
      but it can also be set per path.

    paths                       list
      Paths are what make this application work. Add at least one. Each path should
      have either repo or redir set. Or both. Each path has the following optional
      attributes:

      cache_age
        See explanation above.

      redir_paths
        See explanation above.

      repo
        URL to the repo for the vanity path.

      redir
        URL to redirect this vanity path to. May be used with repo vanity URL.

      display
        This is used as the display line in the go-import html meta tag.
        Determined automatically. Set it if you're going extra custom.

      vcs
        Supported VCSs: git hg bzr svn
        Set this if you get an error that it cannot be determined automatically.

      wildcard
        Allows redirecting all sub paths as repo paths. Set to true to enable the feature.

AUTHOR
---
*   GoogleCloudPlatform - 2017-2018
*   David Newhall II    - 2019

LOCATION
---
*   Turbo Vanity URLs: [https://github.com/golift/turbovanityurls](https://github.com/golift/turbovanityurls)

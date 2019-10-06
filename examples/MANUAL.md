turbovanityurls(1) -- Provide vanity go import paths
===

SYNOPSIS
---
`turbovanityurls -c /etc/turbovanityurls/config.yaml`

This daemon provides a web server that creates go "vanity" import paths.

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

    host                        required
      Used as the import path host. This must be set.

    description
      Displayed as a description paragraph on the index page.

    logo_url
      If this is set, the logo is displayed in the index and package templates.
      Set this to a URI or URL for an image that is used in an img src tag.

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

    redir_index
      If set, this parameter is used to redirect index page requests. By default
      the index page is displayed from a built-in template. If you would rather
      visitors get forwarded to another page put that URI or URL here.

    redir_404
      If set, this parameter is used to redirect 404 requests. Set this to a URI
      or URL to redirect requests to that resulted in a missing page.

    paths                       list
      Paths are what make this application work. Add at least one. Each path should
      have either repo or redir set. Or both. Each path has the following optional
      attributes:

      links
        Each package can have a list of resource links. Each link has a url and
        a title attribute.

      description
        If a description is provided, it is displayed on the package page.
        HTML is OK.

      image_url
        If parameter is provided, it will be displayed at the top of the
        package's page. This must be a URI or URL to an image (png, jpg, etc).

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

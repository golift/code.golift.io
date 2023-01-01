turbo-vanityURLs - golift.io source
---

This is the source that runs [https://golift.io](https://golift.io).

# Automatic Builds

This can be run anywhere, really. Packages and builds for lots of things are provided.
The repo auto-builds packages for freebsd, linux, binaries for macOS and windows. A homebrew
formula and a Docker image are also available, and easy to use.

# Install

- Linux users can use this script to download and install the latest package for their system.<br>
Linux repository hosting provided by
[![packagecloud](https://docs.golift.io/integrations/packagecloud-full.png "PackageCloud.io")](http://packagecloud.io)<br>
This works on any system with apt or yum. If your system does not use APT or YUM, then download a package from the [Releases](https://github.com/Notifiarr/notifiarr/releases) page.

Install the Go Lift package repo and Turbo Vanity URLs with this command:
```
curl -s https://golift.io/repo.sh | sudo bash -s - turbovanityurls
```

- Docker users can pull directly from the image built on Docker.
```
docker pull golift/turbovanityurls
```

The config file is located at `/etc/turbovanityurls/config.yaml`, pass that path into your container.

- macOS users can try it out using homebrew.
```
brew install golift/mugs/turbovanityurls
```

- App Engine:
Run `glcoud app deploy` after you edit [app.yaml](app.yaml).

# Proxy

I run it in Docker behind [swag](https://docs.linuxserver.io/general/swag) (nginx) using a config like this:

```
server {
  # This is the turbovanityurls container.
  set $server http://golift:8080;
  server_name golift.io www.golift.io code.golift.io;

  listen 443 ssl http2;

  location @proxy {
    include  /config/nginx/proxy.conf;
    proxy_pass $server$request_uri;
  }

  location / {
    # This points to the 'static' folder in this repo.
    root /config/www/golift.io;
    try_files $uri @proxy;
  }
}
```

- FreeBSD users can find a package on the [Releases](https://github.com/golift/turbovanityurls/releases) page.

# Changes

Differences from [https://github.com/GoogleCloudPlatform/govanityurls](https://github.com/GoogleCloudPlatform/govanityurls):

-   **Wildcard Support** Example: You can point a path (even /) to a github user/org. [#25](https://github.com/GoogleCloudPlatform/govanityurls/pull/25)
-   App Engine Go 1.12. `go112` [#29](https://github.com/GoogleCloudPlatform/govanityurls/pull/29)
-   App Engine Go 1.15+ `go115`.
-   Moved Templates to their own file.
-   Cleaned up templates. Add some css, a little better formatting.
-   Pass entire `PathConfig` into templates.
-   Exported most of the struct members to make them usable by `yaml` and `template` packages.
-   Reused structs for unmarshalling and passing into templates.
-   Converted `PathConfig` to a pointer; can be accessed as a map or a slice now.
-   Embedded structs for better inheritance model.
-   Set `max_age` per path instead of global-only.
-   Added `-l` (listen), `-t` (timeout), and `-c` (config) flags. [#20](https://github.com/GoogleCloudPlatform/govanityurls/pull/20)
-   Root path repos work now. [#23](https://github.com/GoogleCloudPlatform/govanityurls/pull/23)
-   Better auto-detection for repo type. [#26](https://github.com/GoogleCloudPlatform/govanityurls/pull/26) and [#27](https://github.com/GoogleCloudPlatform/govanityurls/pull/27)

## New Features
-   See the [new manual](examples/MANUAL.md), and the [example config file](examples/config.yaml.example).
-   Path redirects. Issue 302s for specific paths.
    -   Useful for redirecting to download links on GitHub.
-   More customization for index and package pages.
-   Configurable descriptions and logos.
-   Better CSS/HTML templates.

## Other
Incorporated a badge package for data collection and return.
In other words this app can collect data from "things"
(like the public grafana api) and store that data for later requests.
I use this to populate badge/shield data for things like "grafana
dashboard download counter" - [https://github.com/golift/badgedata](https://github.com/golift/badgedata). It's [3 lines of code](https://github.com/golift/turbovanityurls/commit/89451a0a783b9c1991313c0a5cc6e70e9c023e14#diff-7ddfb3e035b42cd70649cc33393fe32c) you can pull out real easy. You can also disable badgedata in the config file.

# Integrations

The following fine folks are providing their services, completely free! These service
integrations are used for things like storage, building, compiling, distribution and
documentation support. This project succeeds because of them. Thank you!

<p style="text-align: center;">
<a title="PackageCloud" alt="PackageCloud" href="https://packagecloud.io"><img src="https://docs.golift.io/integrations/packagecloud.png"/></a>
<a title="GitHub" alt="GitHub" href="https://GitHub.com"><img src="https://docs.golift.io/integrations/octocat.png"/></a>
<a title="Docker Cloud" alt="Docker" href="https://cloud.docker.com"><img src="https://docs.golift.io/integrations/docker.png"/></a>
<a title="Homebrew" alt="Homebrew" href="https://brew.sh"><img src="https://docs.golift.io/integrations/homebrew.png"/></a>
<a title="Go Lift" alt="Go Lift" href="https://golift.io"><img src="https://docs.golift.io/integrations/golift.png"/></a>
</p>

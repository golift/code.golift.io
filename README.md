code.golift.io source
---

This is the source that runs [https://code.golift.io](https://code.golift.io).

#### Fixes from [https://github.com/GoogleCloudPlatform/govanityurls](https://github.com/GoogleCloudPlatform/govanityurls)

-   **Wildcard Support** Example: You can point a path (even /) to a github user/org. [#25](https://github.com/GoogleCloudPlatform/govanityurls/pull/25)
-   App Engine Go 1.12. `go112` [#29](https://github.com/GoogleCloudPlatform/govanityurls/pull/29)
-   Moved Templates to their own file.
-   Cleaned up templates. Add some css, a little better formatting.
-   Pass entire `PathConfig` into templates.
-   Exported most of the struct members to make them usable by `yaml` and `template` packages.
-   Reused structs for unmarshalling and passing into templates.
-   Converted `PathConfig` to a pointer; can be accessed as a map or a slice now.
-   Embedded structs for better inheritance model.
-   Set `max_age` per path instead of global-only.
-   Added `-listen` and `-config` flags. [#20](https://github.com/GoogleCloudPlatform/govanityurls/pull/20)
-   Root path repos work now. [#23](https://github.com/GoogleCloudPlatform/govanityurls/pull/23)
-   Better auto-detection for repo type. [#26](https://github.com/GoogleCloudPlatform/govanityurls/pull/26) and [#27](https://github.com/GoogleCloudPlatform/govanityurls/pull/27)

#### New Features
-   Path redirects. Issue 302s for specific paths.
    -   Useful for redirecting to download links on GitHub.
-   More customization for index page.

#### Other
Incorporated a badge package for data collection and return.
In other words this app can collect data from "things"
(like the public grafana api) and store that data for later requests.
I use this to populate badge/shield data for things like "grafana
dashboard download counter" - [https://github.com/golift/badgedata](https://github.com/golift/badgedata). It's [3 lines of code](https://github.com/golift/code.golift.io/commit/89451a0a783b9c1991313c0a5cc6e70e9c023e14#diff-7ddfb3e035b42cd70649cc33393fe32c) you can pull out real easy. 

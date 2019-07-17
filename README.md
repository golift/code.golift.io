code.golift.io source
---

This is the source that runs [https://code.golift.io](https://code.golift.io).

#### Fixes from [https://github.com/GoogleCloudPlatform/govanityurls](https://github.com/GoogleCloudPlatform/govanityurls)

-   App Engine Go 1.12. `go112`
-   Moved Templates to their own file.
-   Cleaned up templates. Add some css, a little better formatting.
-   Pass entire `PathConfig` into templates.
-   Exported most of the struct members to make them usable by `yaml` and `template` packages.
-   Reused structs for unmarshalling and passing into templates.
-   Converted `PathConfig` to a pointer; can be accessed as a map or a slice now.
-   Embedded structs for better inheritance model.
-   Set `max_age` per path instead of global-only.

#### New Features
-   Path redirects. Issue 302s for specific paths.
    -   Useful for redirecting to download links on GitHub.

#### Other
Incorporated a badge package for data collection and return.
In other words this app can collect data from "things"
(like the public grafana api) and store that data for later requests.
I use this to populate badge/shield data for things like "grafana
dashboard download counter" - [https://github.com/golift/badgedata](https://github.com/golift/badgedata)

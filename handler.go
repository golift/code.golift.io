// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// govanityurls serves Go vanity URLs.
package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// Handler contains all the running data for our web server.
type Handler struct {
	*Config
	PathConfigs
}

// vcsPrefixMap provides defaults for VCS type if it's not provided.
// The list of strings is used in strings.HasPrefix().
var vcsPrefixMap = map[string][]string{
	"git": {"https://git", "https://bitbucket"},
	"bzr": {"https://bazaar"},
	"hg":  {"https://hg.", "https://mercurial"},
	"svn": {"https://svn."},
}

// PathConfigs contains our list of configured routing-paths.
type PathConfigs []*PathConfig

// PathConfig is the configuration for a single routing path.
type PathConfig struct {
	Path     string  `yaml:"-"`
	CacheAge *uint64 `yaml:"cache_max_age,omitempty"`
	ImageURL string  `yaml:"image_url,omitempty"`
	Links    []struct {
		Title string `yaml:"title,omitempty"`
		URL   string `yaml:"url,omitempty"`
	} `yaml:"links,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	RedirPaths   []string `yaml:"redir_paths,omitempty"`
	Repo         string   `yaml:"repo,omitempty"`
	Redir        string   `yaml:"redir,omitempty"`
	Display      string   `yaml:"display,omitempty"`
	VCS          string   `yaml:"vcs,omitempty"`
	Wildcard     bool     `yaml:"wildcard,omitempty"`
	cacheControl string
}

// PathReq is returned by find() with a non-nil PathConfig
// when a request has been matched to a path.
// Host, LogoURL, and IndexTitle come unset.
// This struct is passed into the vanity template.
type PathReq struct {
	Host       string
	Subpath    string
	IndexTitle string
	LogoURL    string
	*PathConfig
}

func (c *Config) newHandler() (*Handler, error) {
	h := &Handler{Config: c}

	if c.Host == "" {
		return nil, fmt.Errorf("must provide host value in config")
	}

	for p := range h.Paths {
		h.Paths[p].Path = p

		if len(h.Paths[p].RedirPaths) < 1 {
			// was not provided, pass in global value.
			h.Paths[p].RedirPaths = h.RedirPaths
		}

		h.Paths[p].setRepoCacheControl(h.CacheAge)

		if err := h.Paths[p].setRepoVCS(); err != nil {
			return nil, err
		}

		h.PathConfigs = append(h.PathConfigs, h.Paths[p])
	}

	sort.Sort(h.PathConfigs)

	return h, nil
}

func (p *PathConfig) setRepoCacheControl(globalCC *uint64) {
	switch {
	case p.CacheAge != nil:
		p.cacheControl = fmt.Sprintf("public, max-age=%d", *p.CacheAge)
	case globalCC != nil:
		p.cacheControl = fmt.Sprintf("public, max-age=%d", *globalCC)
	default:
		p.cacheControl = fmt.Sprintf("public, max-age=86400") // 24 hours (in seconds)
	}
}

// setRepoVCS makes sure the provides VCS type is supported,
// or sets it automatically based on the repo's prefix.
func (p *PathConfig) setRepoVCS() (err error) {
	// Check and set VCS type.
	switch {
	case p.Repo == "" && p.Redir != "":
		// Redirect-only can go anywhere.
	case p.VCS == "github" || p.VCS == "gitlab" || p.VCS == "bitbucket":
		p.VCS = "git"
	case p.VCS == "":
		// Try to figure it out.
		p.VCS, err = findRepoVCS(p.Repo)

		return err
	default:
		// Already filled in, make sure it's supported.
		if _, ok := vcsPrefixMap[p.VCS]; !ok {
			return fmt.Errorf("configuration for %v: unknown VCS %s", p, p.VCS)
		}
	}

	return nil
}

// findRepoVCS checks the vcsMapList for a supported vcs type based on a repo's prefix.
func findRepoVCS(repo string) (string, error) {
	for vcs, prefixList := range vcsPrefixMap {
		for _, prefix := range prefixList {
			if strings.HasPrefix(repo, prefix) {
				return vcs, nil
			}
		}
	}

	return "", fmt.Errorf("cannot infer VCS from %s", repo)
}

// NotFound redirects 404 requests if a redirect URL is set.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	if h.Redir404 != "" {
		http.Redirect(w, r, h.Redir404, http.StatusFound)

		return
	}

	http.NotFound(w, r)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch pc := h.PathConfigs.find(r.URL.Path); {
	case pc.PathConfig == nil && r.URL.Path != "/":
		// Unknown URI
		h.NotFound(w, r)
	case pc.PathConfig == nil && h.RedirIndex != "":
		// Index page, but redirect is present.
		http.Redirect(w, r, h.RedirIndex, http.StatusFound)
	case pc.PathConfig == nil:
		// Index page template.
		templ := indexTmpl.Funcs(funcMap)
		if err := templ.Execute(w, &h.Config); err != nil {
			http.Error(w, "cannot render the page", http.StatusInternalServerError)
		}
	case pc.RedirectablePath():
		// Redirect for file downloads.
		redirTo := pc.Redir + strings.TrimPrefix(r.URL.Path, pc.Path)
		http.Redirect(w, r, redirTo, http.StatusFound)
	case pc.Repo == "":
		// Repo is not set and no paths to redirect, so we're done.
		h.NotFound(w, r)
	default:
		// Create a vanity redirect page.
		w.Header().Set("Cache-Control", pc.cacheControl)
		pc.Host = h.Host
		pc.IndexTitle = h.Title
		pc.LogoURL = h.LogoURL
		templ := vanityTmpl.Funcs(funcMap)

		if r.URL.Query().Get("go-get") == "1" {
			// Use a smaller html page if this is a go-get request.
			templ = gogetTmpl.Funcs(funcMap)
		}

		if err := templ.Execute(w, &pc); err != nil {
			http.Error(w, "cannot render the page", http.StatusInternalServerError)
		}
	}
}

// RedirectablePath checks if a string exists in a list of strings.
// Used to determine if a sub path should be redirected or not.
// Not used for normal vanity URLs, only used for `redir`.
func (p *PathReq) RedirectablePath() bool {
	if p.Redir == "" {
		return false
	}

	for _, s := range p.RedirPaths {
		if strings.Contains(p.Subpath, s) {
			return true
		}
	}

	return false
}

// ImportPath is used in the template to generate the import path.
func (p *PathReq) ImportPath() string {
	path := p.Path

	if p.Wildcard {
		path += strings.Split(p.Subpath, "/")[0]
	}

	return strings.TrimSuffix(path, "/")
}

// RepoPath is used in the template to generate the repo path.
func (p *PathReq) RepoPath() string {
	repo := p.Repo

	if p.Wildcard {
		repo += strings.Split(p.Subpath, "/")[0]
	}

	return repo
}

// Title is used in the template to generate the package title (name).
func (p *PathReq) Title() string {
	s := strings.Split(p.ImportPath(), "/")

	return s[len(s)-1]
}

// SourcePath is used in the template to generate the source path.
func (p *PathReq) SourcePath() string {
	if p.Display != "" {
		return p.Host + p.ImportPath() + " " + p.Display
	}

	template := "%v%v %v %v/tree/master{/dir} %v/blob/master{/dir}/{file}#L{line}"

	if strings.HasPrefix(p.Repo, "https://bitbucket.org") {
		template = "%v%v %v %v/src/default{/dir} %v/src/default{/dir}/{file}#{file}-{line}"
	}

	path := p.ImportPath()
	repo := p.RepoPath()

	// github, gitlab, git, svn, hg, bzr - may need more tweaking for some of these.
	return fmt.Sprintf(template, p.Host, path, repo, repo, repo)
}

// Len is a sort.Sort interface method.
func (pset PathConfigs) Len() int {
	return len(pset)
}

// Less is a sort.Sort interface method.
func (pset PathConfigs) Less(i, j int) bool {
	return pset[i].Path < pset[j].Path
}

// Swap is a sort.Sort interface method.
func (pset PathConfigs) Swap(i, j int) {
	pset[i], pset[j] = pset[j], pset[i]
}

func (pset PathConfigs) find(path string) PathReq {
	// Fast path with binary search to retrieve exact matches
	// e.g. given pset ["/", "/abc", "/xyz"], path "/def" won't match.
	i := sort.Search(len(pset), func(i int) bool {
		return pset[i].Path >= path
	})

	if i < len(pset) && pset[i].Path == path {
		// We have an exact match to a configured path.
		return PathReq{PathConfig: pset[i]}
	}

	// This attempts to match /some/path/here but not /some/pathhere
	if i > 0 && strings.HasPrefix(path, pset[i-1].Path+"/") {
		// We have a partial match with a subpath!
		return PathReq{
			PathConfig: pset[i-1],
			Subpath:    path[len(pset[i-1].Path)+1:],
		}
	}

	// Slow path, now looking for the longest prefix/shortest subpath i.e.
	// e.g. given pset ["/", "/abc", "/abc/def", "/xyz"]
	//  * query "/abc/foo" returns "/abc" with a subpath of "foo"
	//  * query "/x" returns "/" with a subpath of "x"
	lenShortestSubpath := len(path)
	var p PathReq

	// After binary search with the >= lexicographic comparison,
	// nothing greater than i will be a prefix of path.
	max := i
	for i := 0; i < max; i++ {
		if len(pset[i].Path) >= len(path) {
			// We previously didn't find the request path by search, so any
			// configured path with equal or greater length is NOT a match.
			continue
		}

		sSubpath := strings.TrimPrefix(path, pset[i].Path)
		if len(sSubpath) < lenShortestSubpath {
			// We get into this if statement only if TrimPrefix trimmed something.
			// Then we try the next path, and check to see if what's left after we
			// trimmed the configured path off is shorter than this one. /x is better than /xyz.
			p.Subpath = sSubpath
			lenShortestSubpath = len(sSubpath)
			p.PathConfig = pset[i]
		}
	}

	return p
}

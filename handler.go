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

	"google.golang.org/appengine"
	yaml "gopkg.in/yaml.v2"
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

// Config contains the config file data.
type Config struct {
	Title string `yaml:"title,omitempty"`
	Host  string `yaml:"host,omitempty"`
	Links []struct {
		Title string `yaml:"title,omitempty"`
		URL   string `yaml:"url,omitempty"`
	} `yaml:"links,omitempty"`
	CacheAge   *uint64                `yaml:"cache_max_age,omitempty"`
	Paths      map[string]*PathConfig `yaml:"paths,omitempty"`
	RedirPaths []string               `yaml:"redir_paths,omitempty"`
	Src        string                 `yaml:"src,omitempty"`
}

// PathConfigs contains our list of configured routing-paths.
type PathConfigs []*PathConfig

// PathConfig is the configuration for a single routing path.
type PathConfig struct {
	Path         string            `yaml:"-"`
	CacheAge     *uint64           `yaml:"cache_max_age,omitempty"`
	RedirPaths   []string          `yaml:"redir_paths,omitempty"`
	Repo         string            `yaml:"repo,omitempty"`
	Redir        string            `yaml:"redir,omitempty"`
	Display      string            `yaml:"display,omitempty"`
	VCS          string            `yaml:"vcs,omitempty"`
	Wildcard     bool              `yaml:"wildcard,omitempty"`
	Tags         map[string]string `yaml:"tags,omitempty"`
	cacheControl string
}

// PathReq is returned by find() with a non-nil PathConfig
// when a request has been matched to a path. Host comes unset.
// This struct is passed into the vanity template.
type PathReq struct {
	Host    string
	Subpath string
	Tag     string
	*PathConfig
}

func newHandler(configData []byte) (*Handler, error) {
	h := &Handler{Config: &Config{Paths: make(map[string]*PathConfig)}}
	if err := yaml.Unmarshal(configData, h.Config); err != nil {
		return nil, err
	}
	if h.Title == "" {
		h.Title = h.Host
	}
	cacheControl := fmt.Sprintf("public, max-age=86400") // 24 hours (in seconds)
	if h.CacheAge != nil {
		cacheControl = fmt.Sprintf("public, max-age=%d", *h.CacheAge)
	}
	for p := range h.Paths {
		h.Paths[p].Path = p
		if len(h.Paths[p].RedirPaths) < 1 {
			// was not provided, pass in global value.
			h.Paths[p].RedirPaths = h.RedirPaths
		}
		h.Paths[p].setRepoCacheControl(cacheControl)
		if err := h.Paths[p].setRepoVCS(); err != nil {
			return nil, err
		}
		h.PathConfigs = append(h.PathConfigs, h.Paths[p])
	}
	sort.Sort(h.PathConfigs)
	return h, nil
}

func (p *PathConfig) setRepoCacheControl(cc string) {
	p.cacheControl = cc
	if p.CacheAge != nil {
		// provided, override global value.
		p.cacheControl = fmt.Sprintf("public, max-age=%d", *p.CacheAge)
	}
}

// setRepoVCS makes sure the provides VCS type is supported,
// or sets it automatically based on the repo's prefix.
func (p *PathConfig) setRepoVCS() error {
	// Check and set VCS type.
	switch {
	case p.Repo == "" && p.Redir != "":
		// Redirect-only can go anywhere.
	case p.VCS == "github" || p.VCS == "gitlab" || p.VCS == "bitbucket":
		p.VCS = "git"
	case p.VCS == "":
		// Try to figure it out.
		var err error
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pc := h.PathConfigs.find(r.URL.Path)
	if pc.PathConfig == nil {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		} else if err := indexTmpl.Execute(w, &h.Config); err != nil {
			http.Error(w, "cannot render the page", http.StatusInternalServerError)
		}
		return
	}
	if pc.RedirectablePath() {
		// Redirect for file downloads.
		redirTo := pc.Redir + strings.TrimPrefix(r.URL.Path, pc.Path)
		http.Redirect(w, r, redirTo, http.StatusFound)
		return
	}
	if pc.Repo == "" {
		// Repo is not set and no paths to redirect, so we're done.
		http.NotFound(w, r)
		return
	}

	// Create a vanity redirect page.
	w.Header().Set("Cache-Control", pc.cacheControl)
	pc.Host = h.Hostname(r)
	if err := vanityTmpl.Execute(w, &pc); err != nil {
		http.Error(w, "cannot render the page", http.StatusInternalServerError)
	}
}

// Hostname returns the appropriate Host header for this request.
func (h *Handler) Hostname(r *http.Request) string {
	if h.Host != "" {
		return h.Host
	}
	appHost := appengine.DefaultVersionHostname(appengine.NewContext(r))
	if appHost == "" {
		return r.Host
	}
	return appHost
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
	_, stag := p.GetTag()
	repo := p.Repo
	if p.Wildcard {
		sub := strings.Split(p.Subpath, "/")[0]
		path += sub
		repo += sub
	}
	return fmt.Sprintf("%v%v%v %v %v", p.Host, strings.TrimSuffix(path, "/"), stag, p.VCS, repo)
}

// SourcePath is used in the template to generate the source path.
func (p *PathReq) SourcePath() string {
	if p.Display != "" {
		return p.Host + strings.TrimSuffix(p.Path, "/") + " " + p.Display
	}
	tag, stag := p.GetTag()
	template := "%v%v%v %v %v/tree/%v{/dir} %v/blob/%v{/dir}/{file}#L{line}"
	if strings.HasPrefix(p.Repo, "https://bitbucket.org") {
		if tag == "master" {
			tag = "default"
		}
		template = "%v%v%v %v %v/src/%v{/dir} %v/src/%v{/dir}/{file}#{file}-{line}"
	}
	path := p.Path
	repo := p.Repo
	if p.Wildcard {
		sub := strings.Split(p.Subpath, "/")[0]
		path += sub
		repo += sub
	}
	// DO not include the middle rep path if we have a source tag.
	mrepo := repo
	if stag != "" {
		mrepo = "_"
	}
	// github, gitlab, git, svn, hg, bzr - may need more tweaking for some of these.
	return fmt.Sprintf(template, p.Host, strings.TrimSuffix(path, "/"), stag, mrepo, repo, tag, repo, tag)
}

// GetTag is a helper to return the tag and branch for a request.
func (p *PathReq) GetTag() (string, string) {
	tag := "master"
	stag := ""
	if t, ok := p.Tags[p.Tag]; ok {
		tag = t
		stag = "." + p.Tag
	}
	return tag, stag
}

// GoDocPath is used in the template to generate the GoDoc path.
func (p *PathReq) GoDocPath() string {
	if p.Wildcard {
		return p.Host + "/" + p.Subpath
	}
	return p.Host + p.Path + "/" + p.Subpath
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
	var p PathReq
	sp := strings.Split(path, "/")             // get last path element.
	sp = strings.SplitN(sp[len(sp)-1], ".", 2) // check if it has a dot.
	if len(sp) == 2 {                          // dot at the end of a path refers to a tag.
		p.Tag = sp[1] // Save the tag.
	}
	// Fast path with binary search to retrieve exact matches
	// e.g. given pset ["/", "/abc", "/xyz"], path "/def" won't match.
	i := sort.Search(len(pset), func(i int) bool {
		return pset[i].Path >= path
	})
	if i < len(pset) && pset[i].Path == path {
		// We have an exact match to a configured path.
		p.PathConfig = pset[i]
		return p
	}
	// This attempts to match /some/path/here but not /some/pathhere
	if i > 0 && strings.HasPrefix(path, pset[i-1].Path+"/") {
		// We have a partial match with a subpath!
		p.PathConfig = pset[i-1]
		p.Subpath = path[len(pset[i-1].Path)+1:]
		return p
	}

	// Slow path, now looking for the longest prefix/shortest subpath i.e.
	// e.g. given pset ["/", "/abc", "/abc/def", "/xyz"]
	//  * query "/abc/foo" returns "/abc" with a subpath of "foo"
	//  * query "/x" returns "/" with a subpath of "x"
	lenShortestSubpath := len(path)

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

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

// nolint:forbidigo
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golift.io/badgedata"
	_ "golift.io/badgedata/grafana"
	yaml "gopkg.in/yaml.v3"
)

// Version is injected at build time.
var Version = "development" //nolint:gochecknoglobals

// Flags are the CLI flags.
type Flags struct {
	listenAddr string
	configPath string
	showVer    bool
}

// Config contains the config file data.
type Config struct {
	Title       string `yaml:"title,omitempty"`
	Host        string `yaml:"host,omitempty"`
	Description string `yaml:"description,omitempty"`
	LogoURL     string `yaml:"logo_url,omitempty"`
	Links       []struct {
		Title string `yaml:"title,omitempty"`
		URL   string `yaml:"url,omitempty"`
	} `yaml:"links,omitempty"`
	CacheAge   *uint64                `yaml:"cache_max_age,omitempty"`
	BDPath     string                 `yaml:"bd_path,omitempty"`
	Paths      map[string]*PathConfig `yaml:"paths,omitempty"`
	RedirPaths []string               `yaml:"redir_paths,omitempty"`
	Src        string                 `yaml:"src,omitempty"`
	RedirIndex string                 `yaml:"redir_index,omitempty"`
	Redir404   string                 `yaml:"redir_404,omitempty"`
}

func parseFlags(args []string) *Flags {
	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f := &Flags{listenAddr: ":" + os.Getenv("PORT")}

	if f.listenAddr == ":" {
		f.listenAddr = ":8080"
	}

	flag.StringVar(&f.listenAddr, "l", f.listenAddr, "HTTP server listen address")
	flag.StringVar(&f.configPath, "c", DefaultConfFile, "config file path")
	flag.BoolVar(&f.showVer, "v", false, "show version and exit")

	flag.Usage = func() {
		fmt.Println("Usage: turbovanityurls [-c <config-file>] [-l <listen-address>]")
		flag.PrintDefaults()
	}

	_ = flag.Parse(args)

	return f
}

func main() {
	flags := parseFlags(os.Args[1:])
	if flags.showVer {
		fmt.Printf("%v v%v\n", "turbovanityurls", Version)
		os.Exit(0)
	}

	config, err := parseConfig(flags.configPath)
	if err != nil {
		log.Fatal(err)
	}

	vanityHandler, err := config.newHandler()
	if err != nil {
		log.Fatal(err)
	}

	if config.BDPath != "" {
		http.Handle(config.BDPath, badgedata.Handler())
	}

	http.Handle("/", vanityHandler)

	if strings.HasPrefix(flags.listenAddr, ":") {
		// A message so you know when it's started; a clickable link for dev'ing.
		log.Println("Listening at http://127.0.0.1" + flags.listenAddr)
	}

	if err := http.ListenAndServe(flags.listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(configPath string) (*Config, error) {
	c := &Config{Paths: make(map[string]*PathConfig)}

	if _, err := os.Stat(configPath); os.IsNotExist(err) && configPath == DefaultConfFile {
		log.Printf("Default Config File Not Found: %s - trying ./config.yaml", configPath)
		configPath = "config.yaml"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, fmt.Errorf("unmarshaling config file: %w", err)
	}

	if c.Title == "" {
		c.Title = c.Host
	}

	if c.BDPath != "" && !strings.HasSuffix(c.BDPath, "/") {
		c.BDPath += "/"
	}

	return c, nil
}

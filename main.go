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
	"golift.io/turbovanityurls/pkg/handler"
	yaml "gopkg.in/yaml.v3"
)

// Version is injected at build time.
var Version = "development" //nolint:gochecknoglobals

// Flags are the CLI flags.
type Flags struct {
	ListenAddr string
	ConfigPath string
	ShowVer    bool
}

type Config struct {
	*handler.Config `yaml:",inline"`
	BDPath          string `yaml:"bd_path,omitempty"`
}

func ParseFlags(args []string) *Flags {
	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f := &Flags{ListenAddr: ":" + os.Getenv("PORT")}

	if f.ListenAddr == ":" {
		f.ListenAddr = ":8080"
	}

	flag.StringVar(&f.ListenAddr, "l", f.ListenAddr, "HTTP server listen address")
	flag.StringVar(&f.ConfigPath, "c", DefaultConfFile, "config file path")
	flag.BoolVar(&f.ShowVer, "v", false, "show version and exit")

	flag.Usage = func() {
		fmt.Println("Usage: turbovanityurls [-c <config-file>] [-l <listen-address>]")
		flag.PrintDefaults()
	}

	_ = flag.Parse(args)

	return f
}

func main() {
	flags := ParseFlags(os.Args[1:])
	if flags.ShowVer {
		fmt.Printf("%v v%v\n", "turbovanityurls", Version)
		os.Exit(0)
	}

	if err := Setup(flags); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(flags.ListenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func Setup(flags *Flags) error {
	config := &Config{}
	if err := config.ParseConfig(flags.ConfigPath); err != nil {
		return err
	}

	vanityHandler, err := handler.New(config.Config)
	if err != nil {
		return fmt.Errorf("config file: %w", err)
	}

	if config.BDPath != "" {
		http.Handle(config.BDPath, badgedata.Handler())
	}

	http.Handle("/", vanityHandler)

	if strings.HasPrefix(flags.ListenAddr, ":") {
		// A message so you know when it's started; a clickable link for dev'ing.
		log.Println("Listening at http://127.0.0.1" + flags.ListenAddr)
	}

	return nil
}

func (c *Config) ParseConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) && configPath == DefaultConfFile {
		log.Printf("Default Config File Not Found: %s - trying ./config.yaml", configPath)
		configPath = "config.yaml"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("unmarshaling config file: %w", err)
	}

	if c.Title == "" {
		c.Title = c.Host
	}

	if c.BDPath != "" && !strings.HasSuffix(c.BDPath, "/") {
		c.BDPath += "/"
	}

	return nil
}

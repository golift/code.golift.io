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
)

// Version is injected at build time.
var Version = "development"

// Flags are the CLI flags.
type Flags struct {
	listenAddr string
	configPath string
	showVer    bool
}

func main() {
	flags := parseFlags()
	if flags.showVer {
		fmt.Printf("%v v%v\n", "turbovanityurls", Version)
		os.Exit(0)
	}
	vanity, err := ioutil.ReadFile(flags.configPath)
	if err != nil {
		log.Fatal(err)
	}
	vanityHandler, err := newHandler(vanity)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/bd/", badgedata.Handler())
	http.Handle("/", vanityHandler)
	if strings.HasPrefix(flags.listenAddr, ":") {
		// A message so you know when it's started; a clickable link for dev'ing.
		log.Println("Listening at http://127.0.0.1" + flags.listenAddr)
	}
	if err := http.ListenAndServe(flags.listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() *Flags {
	f := &Flags{listenAddr: ":" + os.Getenv("PORT")}
	if f.listenAddr == ":" {
		f.listenAddr = ":8080"
	}
	flag.StringVar(&f.listenAddr, "l", f.listenAddr, "HTTP server listen address")
	flag.StringVar(&f.configPath, "c", "./config.yaml", "config file path")
	flag.BoolVar(&f.showVer, "v", false, "show version and exit")
	flag.Usage = func() {
		fmt.Println("Usage: turbovanityurls [-c <config-file>] [-l <listen-address>]")
		flag.PrintDefaults()
	}
	flag.Parse()
	return f
}

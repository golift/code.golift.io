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

	"code.golift.io/badgedata"
	_ "code.golift.io/badgedata/grafana"
)

// Version is injected at build time.
var Version = "development"

func main() {
	listenAddr, configPath := parseFlags()
	vanity, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	vanityHandler, err := newHandler(vanity)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/bd/", badgedata.Handler())
	http.Handle("/", vanityHandler)
	if strings.HasPrefix(listenAddr, ":") {
		// A message so you know when it's started; a clickable link for dev'ing.
		log.Println("Listening at http://127.0.0.1" + listenAddr)
	}
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() (string, string) {
	listenAddr := ":" + os.Getenv("PORT")
	if listenAddr == ":" {
		listenAddr = ":8080"
	}
	flag.StringVar(&listenAddr, "l", listenAddr, "HTTP server listen address")
	configPath := flag.String("c", "./config.yaml", "config file path")
	showVer := flag.Bool("v", false, "show version and exit")
	flag.Usage = func() {
		fmt.Println("Usage: turbovanityurls [-c <config-file>] [-l <listen-address>]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *showVer {
		// move this into main.
		fmt.Printf("%v v%v\n", "turbovanityurls", Version)
		os.Exit(0)
	}
	return listenAddr, *configPath
}

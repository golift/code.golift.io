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

	"code.golift.io/badgedata"
	_ "code.golift.io/badgedata/grafana"
)

func main() {
	listenAddr := os.Getenv("PORT")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	flag.StringVar(&listenAddr, "listen", listenAddr, "HTTP server listen address")
	configPath := flag.String("config", "./config.yaml", "config file path")
	flag.Usage = func() {
		fmt.Println("Usage: govanityurls [-config <config-file>] [-listen <listen-address>]")
		flag.PrintDefaults()
	}
	flag.Parse()

	vanity, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	vanityHandler, err := newHandler(vanity)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/bd/", badgedata.Handler())
	http.Handle("/", vanityHandler)
	// msg is only used to print a message. Useful to know when the app has
	// finished starting and provides a clickable link to get right to it.
	msg := listenAddr
	if msg[0] == ':' {
		msg = "127.0.0.1" + msg
	}
	log.Println("Listening at http://" + msg)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

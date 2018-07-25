// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/centrifugal/gocent"
)

var addr = flag.String("addr", ":8080", "http service address")

var config = struct {
	Secret string `json:"secret"`
}{}

func main() {
	flag.Parse()

	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		panic(err)
	}
	if config.Secret = strings.TrimSpace(config.Secret); config.Secret == "" {
		panic("empty secret!")
	}

	token := auth.GenerateClientToken(config.Secret, "42", "1451991486", "")
	log.Println(token)
	http.Handle("/", http.FileServer(http.Dir("./assets")))

	client := gocent.NewClient("http://localhost:8000", config.Secret, 5*time.Second)
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method != http.MethodPost {
			return
		}
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		success, err := client.Publish("public:news", bodyBytes)
		if err != nil {
			log.Println("could not publish:", err)
		}
		if !success {
			log.Println("success if false")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started on " + *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

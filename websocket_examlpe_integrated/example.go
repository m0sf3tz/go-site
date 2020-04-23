// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

type Person struct {
	FirstName string
	LastName  string
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		p := Person{}
		er := json.Unmarshal(message, &p)
		fmt.Println(er)
		fmt.Println("first: ", p.FirstName, " Last: ", p.LastName)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadFile("./static/button.html")
	w.Write(b)
}

func main() {
	http.HandleFunc("/upgrade", upgrade)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

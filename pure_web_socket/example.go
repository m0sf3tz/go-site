package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

var store = sessions.NewCookieStore([]byte("some key"))

type Person struct {
	FirstName string
	LastName  string
}

func upgrade(w http.ResponseWriter, r *http.Request) {
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

func foo(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "etc")

	session.Values["bar"] = "foo"
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.MaxAge = 0
	session.Save(r, w)
}

func bar(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "etc")
	val := session.Values["bar"]
	s, _ := val.(string)
	fmt.Println(s)
}

func main() {
	http.HandleFunc("/upgrade", upgrade)
	http.HandleFunc("/", home)
	http.HandleFunc("/foo", foo)
	http.HandleFunc("/bar", bar)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	//server.ListenAndServe()
	server.ListenAndServeTLS("cert.pem", "key.pem")

}

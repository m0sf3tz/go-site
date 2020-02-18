package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func login(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("microsoft.html")
	t.Execute(w, "null")
}

func main() {

	fileServer := http.FileServer(http.Dir("microsoft.html"))
	http.Handle("/", fileServer)

	http.ListenAndServe(":8090", nil)
}

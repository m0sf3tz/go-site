package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func in(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "login good")
}
func bad(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "login bad")
}

func check(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	for key, value := range req.Form {
		fmt.Println("%s = %s\n", key, value)
		if key == "uname" {
			if value[0] == "sam" {
				http.Redirect(w, req, "/good", http.StatusOK)
			}
		}
	}

	http.Redirect(w, req, "/bad", http.StatusOK)
}

func login(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {

		req.ParseForm()
		fmt.Println(req.Form)

		if req.Form["uname"][0] == "sam" && req.Form["psw"][0] == "foo" {
			fmt.Println("Correct paissword")
			http.Redirect(w, req, "/in", http.StatusSeeOther)
		} else {
			fmt.Println("Incorrect password")
		}
	}
	t, _ := template.ParseFiles("microsoft.html")
	t.Execute(w, "null")
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/in", in)
	mux.HandleFunc("/bad", bad)
	mux.HandleFunc("/", login)

	fileServer := http.FileServer(http.Dir("./images"))

	// Use the mux.Handle() function to register the file server as the
	// handler for all URL paths that start with "/images/". For matching
	// paths, we strip the "/image" prefix before the request reaches the file
	// server.
	mux.Handle("/images/", http.StripPrefix("/images", fileServer))

	fmt.Println(http.ListenAndServe(":8090", mux))
}

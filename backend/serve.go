package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
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

// sanitized the file-server not to show the "root"
// at ./static/
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/in", in)
	mux.HandleFunc("/bad", bad)
	mux.HandleFunc("/", login)

	fileServer := http.FileServer(http.Dir("./images"))
	mux.Handle("/images/", http.StripPrefix("/images", neuter(fileServer)))

	fmt.Println(http.ListenAndServe(":8090", mux))
}

package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			log.Println("Get")
			tmp, err := template.ParseFiles(
				"./ui/html/index.page.html",
				"./ui/html/header.partial.html",
				"./ui/html/footer.partial.html",
				"./ui/html/base.layout.html",
			)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println("executing")
			if err := tmp.Execute(w, nil); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println("Not Get")
		}

	})

	// serving static files
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	// Using http.StripPrefix to remove leading "/static" from url.
	// This leaves us with "/" which is the current dir of our file server.
	// If we dont strip, our file server will look in the dir "./ui/static/static"
	// rather than "./ui/static/"
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	log.Println("started the server")
	http.ListenAndServe(":4000", mux)
}

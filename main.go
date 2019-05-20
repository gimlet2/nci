package main

import (
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Application started")

	Server(func(r *mux.Router) {
		Get(r, "/", func(w http.ResponseWriter, r *http.Request) {
			Json(w, "hello")
		})
	})
}




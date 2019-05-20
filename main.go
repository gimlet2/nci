package main

import (
	"log"
	"os"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Application started")
	r := mux.NewRouter()
	get(r, "/", func(w http.ResponseWriter, r *http.Request) {
		Json(w, "hello")
	})

	srv := &http.Server{
		Addr: "0.0.0.0:" + getPort(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	log.Fatal(srv.ListenAndServe())
}

func get(r *mux.Router, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	request(r, pattern, handler, "GET")
}

func post(r *mux.Router, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	request(r, pattern, handler, "POST")
}

func request(r *mux.Router, pattern string, handler func(w http.ResponseWriter, r *http.Request), method string) {
	r.HandleFunc(pattern, handler).Methods(method)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

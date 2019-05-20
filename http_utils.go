package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/gorilla/mux"
)

func Json(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(v)
	w.Write(j)
}

func Text(w http.ResponseWriter, v []byte) {
	w.Header().Set("Content-Type", "text")
	w.Write(v)
}

func Error(w http.ResponseWriter, e error) {
	log.Printf("Failed to get resourse %v", e)
	w.WriteHeader(500)
	w.Write([]byte("Error"))
}

func Server(init func(r *mux.Router)) {
	r := mux.NewRouter()
	init(r)
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

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func Get(r *mux.Router, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	request(r, pattern, handler, "GET")
}

func Post(r *mux.Router, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	request(r, pattern, handler, "POST")
}

func request(r *mux.Router, pattern string, handler func(w http.ResponseWriter, r *http.Request), method string) {
	r.HandleFunc(pattern, handler).Methods(method)
}
package main

import (
	"encoding/json"
	"log"
	"net/http"
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

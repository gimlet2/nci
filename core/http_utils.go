package core

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

type Server interface {
	Get(pattern string, handler func(w http.ResponseWriter, r *http.Request))

	Post(pattern string, handler func(w http.ResponseWriter, r *http.Request))

	Json(w http.ResponseWriter, v interface{})

	Html(w http.ResponseWriter, v string)

	ReadJson(r *http.Request) interface{}

	Text(w http.ResponseWriter, v string)

	Error(w http.ResponseWriter, e error)
}

type serverImpl struct {
	r *mux.Router
}

func (s *serverImpl) Get(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	get(s.r, pattern, handler)
}

func (s *serverImpl) Post(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	post(s.r, pattern, handler)
}

func (s *serverImpl) Json(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(v)
	w.Write(j)
}

func (s *serverImpl) Text(w http.ResponseWriter, v string) {
	w.Header().Set("Content-Type", "text")
	w.Write([]byte(v))
}

func (s *serverImpl) Html(w http.ResponseWriter, v string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(v))
}

func (s *serverImpl) Error(w http.ResponseWriter, e error) {
	log.Printf("Failed to get resourse %v", e)
	w.WriteHeader(500)
	w.Write([]byte("Error"))
}

func (s *serverImpl) ReadJson(r *http.Request) interface{} {
	decoder := json.NewDecoder(r.Body)
	var body interface{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Failed to parse json: %v", err)
		return nil
	}
	return body
}

func ServerSetup(init func(s Server)) {
	// autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(os.Getenv("hostname")),
	// 	Cache:      autocert.DirCache(os.Getenv("cacheDir")),
	// }
	r := mux.NewRouter()
	var server Server
	server = &serverImpl{r}
	init(server)
	l := autocert.NewListener(os.Getenv("HOSTNAME"))
	srv := &http.Server{
		// Addr: "0.0.0.0:" + getPort(),
		Addr: ":https",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.

	}
	log.Fatal(srv.Serve(l))

}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
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

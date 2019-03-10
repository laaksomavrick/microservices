package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port = flag.Int("port", 3000, "The http server port")
)

// this should be the entry point to the backend, mapping REST calls to rpc and event publishin

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("ok")
	})
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("catch all")
	})

	log.Printf("Serving on %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

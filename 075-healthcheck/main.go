package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR in /hello, invalid database connection"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

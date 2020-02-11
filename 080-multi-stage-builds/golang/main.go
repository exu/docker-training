package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello!"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

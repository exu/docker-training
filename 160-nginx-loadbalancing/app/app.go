package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Hi from: %s", hostname)
}

func main() {
	http.HandleFunc("/", handler)
	panic(http.ListenAndServe(":8080", nil))
}

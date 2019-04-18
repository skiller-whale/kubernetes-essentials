package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Instance is service instance
type Instance struct {
	// ID is instance identifier
	ID string `json:"id"`
}

// root is root HTTP handler i.e. it handles HTTP queries for "/*"
func root(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %v %v From: %v", r.Method, r.URL.Path, r.RemoteAddr)

		log.Printf("Connecting to: %v", url)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var id Instance
		json.NewDecoder(resp.Body).Decode(&id)
		fmt.Fprintf(w, "<b>Instance: %s</b>", id.ID)
	}
}

func main() {
	// attempt to read listen address
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	// first try to read REMOTE_URL from environment
	url := os.Getenv("REMOTE_URL")
	// if empty, bail
	if url == "" {
		log.Fatal("Empty REMOTE_URL supplied")
	}

	// register HTTP handlers
	http.HandleFunc("/", root(url))

	log.Printf("Starting HTTP server Addr: %v...", addr)

	http.ListenAndServe(":8080", nil)
}

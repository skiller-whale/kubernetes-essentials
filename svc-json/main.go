package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"golang.org/x/net/netutil"
)

// Instance is service instance
type Instance struct {
	// ID is instance identifier
	ID string `json:"id"`
}

// InstanceV2 is service instance with username
type InstanceV2 struct {
	// ID is instance identifier
	ID       string `json:"id"`
	UserName string `json:"username"`
}

// healthz is a liveness probe.
// healthz is a simple HTTP handler which always returns 200
func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readyz is a readiness probe.
// readyz is a HTTP handler which waits for $READY seconds to start responding
func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// root is root HTTP handler i.e. it handles HTTP queries for "/*"
func root(version, delay int, id, path string, isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
      return
    }

		if delay != 0 {
			// if we specified delay, we wait for d ms
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		log.Printf("Received %v %v From: %v", r.Method, r.URL.Path, r.RemoteAddr)

		// if we pass in some path, we need to check if it exists
		// if not we return http.StatusInternalServerError
		if path != "" {
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				log.Printf("Error: %v", err)
				http.Error(w, "Server not ready", http.StatusInternalServerError)
				return

			}
		}

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)

		if version == 2 {
			resp := InstanceV2{
				ID:       id,
				UserName: "Sam Sepiol",
			}

			encoder.Encode(resp)
		} else {
			resp := Instance{
				ID: id,
			}

			encoder.Encode(resp)
		}
	}
}

func main() {
	// attempt to read service listen address
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":80"
	}

	// check ID environment variable
	id := os.Getenv("ID")
	// if empty, generate random UUID
	if id == "" {
		var err error
		id, err = newUUID()
		if err != nil {
			// if we can't generate UUID we bail
			log.Fatalf("Could not start server: %v", err)
		}
	}

	// Version to use
	version := os.Getenv("VERSION")
	if version == "" {
		version = "1"
	}

	v, err := strconv.Atoi(version)
	if err != nil {
		v = 1
	}

	var d int
	// read delay from environment variable
	delay := os.Getenv("DELAY")
	// if defined attempt to convert it to number
	if delay != "" {
		var err error
		d, err = strconv.Atoi(delay)
		if err != nil {
			// if we can't decode we use 500ms
			d = 500
		}
	}

	var r int
	// read readiness time from environment variable
	ready := os.Getenv("READY")
	// if defined convert to number
	if ready != "" {
		var err error
		r, err = strconv.Atoi(ready)
		if err != nil {
			// if we can't decode, then the image should be ready instantly
			r = 0
		}
	}

	// read cache path from environment variable
	path := os.Getenv("CACHE_PATH")

	// atomic readiness flag; no need for channels here
	isReady := &atomic.Value{}
	isReady.Store(false)
	// start a single goroutine to simulate readiness
	go func() {
		log.Printf("Waiting for %v seconds to become ready...", r)
		time.Sleep(time.Duration(r) * time.Second)
		isReady.Store(true)
		log.Printf("Service has become ready now!")
	}()

	// register HTTP handlers
	http.HandleFunc("/", root(v, d, id, path, isReady))
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/readyz", readyz(isReady))

	log.Printf("Starting HTTP server ID: %v Addr: %v...", id, addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}
	defer listener.Close()

	// Restrict to one connection at a time
	listener = netutil.LimitListener(listener, 1)

	log.Fatal(http.Serve(listener, nil))
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

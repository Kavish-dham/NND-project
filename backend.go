package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func createServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	portStr := strconv.Itoa(port)

	// Create a new ServeMux for each server
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", port)
		log.Printf("Server on port %d received request from %s\n", port, r.RemoteAddr)
	})

	server := &http.Server{
		Addr:    ":" + portStr,
		Handler: mux,
	}

	fmt.Printf("Starting server on port %s\n", portStr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server on port %s failed to start: %v", portStr, err)
	}
}

func main() {
	var wg sync.WaitGroup

	for port := 5000; port <= 5003; port++ {
		wg.Add(1)
		go createServer(port, &wg)
	}

	fmt.Println("All servers started. Press Ctrl+C to stop.")
	wg.Wait()
}

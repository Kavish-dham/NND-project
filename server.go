package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Global variables
var (
	backendServers = []string{
		"https://example.com",
		"https://facebook.com",
		"https://instagram.com",
		"https://jsonplaceholder.typicode.com",
		"https://api.publicapis.org",
		"https://dog.ceo/api/breeds/list/all",
	}

	serverResponseTimes = make(map[string]float64)
	mutex               = &sync.Mutex{}
)

func getServerResponse(server string) {
	startTime := time.Now()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	_, err := client.Get(server)

	mutex.Lock()
	defer mutex.Unlock()

	if err != nil {
		serverResponseTimes[server] = -1 // Indicates connection failure
	} else {
		// Convert to milliseconds, with 2 decimal places of precision
		responseTime := float64(time.Since(startTime).Milliseconds())
		serverResponseTimes[server] = responseTime
	}
}

func updateServerStats() {
	for {
		for _, server := range backendServers {
			go getServerResponse(server)
		}
		time.Sleep(3 * time.Second) // Update every 3 seconds
	}
}

func proxyServer(w http.ResponseWriter, r *http.Request) {
	// Choose a random backend server
	server := backendServers[rand.Intn(len(backendServers))]
	http.Redirect(w, r, server, http.StatusTemporaryRedirect)
}

func serverStats(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(serverResponseTimes)
}

func main() {
	// Initialize response times
	for _, server := range backendServers {
		serverResponseTimes[server] = 0
	}

	// Start the background monitoring thread
	go updateServerStats()

	// Set up routes
	http.HandleFunc("/", proxyServer)
	http.HandleFunc("/stats", serverStats)

	// Start the server
	log.Println("Server starting on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

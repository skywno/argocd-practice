package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/time", timeHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /        - Welcome message")
	fmt.Println("  GET /health  - Health check")
	fmt.Println("  GET /time    - Current time")
	
	log.Fatal(http.ListenAndServe(port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Simple Go App</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.container { max-width: 600px; margin: 0 auto; }
				h1 { color: #333; }
				.endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Welcome to Simple Go App!</h1>
				<p>This is a simple Go web application.</p>
				<div class="endpoint">
					<strong>Available endpoints:</strong>
					<ul>
						<li><a href="/health">/health</a> - Health check</li>
						<li><a href="/time">/time</a> - Current time</li>
					</ul>
				</div>
			</div>
		</body>
		</html>
	`)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "healthy", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"time": "%s", "timezone": "UTC"}`, time.Now().UTC().Format(time.RFC3339))
}

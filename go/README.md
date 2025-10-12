# Simple Go App

A simple Go web application with basic HTTP endpoints.

## Features

- HTTP server running on port 8080
- Home page with HTML interface
- Health check endpoint
- Time endpoint
- Clean, responsive web interface

## Running the Application

1. Make sure you have Go installed (version 1.21 or later)
2. Run the application:
   ```bash
   go run main.go
   ```
3. Open your browser and visit `http://localhost:8080`

## Available Endpoints

- `GET /` - Welcome page with HTML interface
- `GET /health` - Health check (returns JSON)
- `GET /time` - Current time (returns JSON)

## Building

To build the application:
```bash
go build -o simple-go-app main.go
```

To run the built binary:
```bash
./simple-go-app
```

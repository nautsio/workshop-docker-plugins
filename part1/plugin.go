package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// Location to the plugin socket
const (
	pluginPath = "/run/docker/plugins/ddd.sock"
)

func main() {
	// Cleanup previous socket if it exists
	if file, _ := os.Stat(pluginPath); file != nil {
		log.Println("Socket already exists, replacing.")
		os.Remove(pluginPath)
	}

	// Create a network socket on the filesystem
	socket, err := net.Listen("unix", pluginPath)
	if err != nil {
		log.Fatalf("unable to listen at %s: %s", pluginPath, err)
	}

	// Create a HTTP request multiplexer (router)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Dutch Docker Day!\n")
	})

	log.Print("Starting Dutch Docker Day volume plugin...")

	// Start the HTTP server
	if err := http.Serve(socket, mux); err != nil {
		log.Fatal("Could not start HTTP server.")
	}
}

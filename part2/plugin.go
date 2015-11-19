package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	pluginPath      = "/run/docker/plugins/ddd.sock"
	versionMimetype = "application/vnd.docker.plugins.v1.1+json"
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

	// Add activation route
	mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) {
		// Docker expects to find a header with the correct mime-type and version
		w.Header().Set("Content-Type", versionMimetype)

		// Create an anonymous struct that is used to encode a JSON response
		// that will let the Docker daemon know that we are capable of handling
		// volume requests.
		//
		// 	{
		//		implements: ["VolumeDriver"]
		//	}
		json.NewEncoder(w).Encode(struct {
			Implements []string
		}{
			[]string{"VolumeDriver"},
		})
	})

	log.Print("Starting Dutch Docker Day volume plugin...")

	// Start the HTTP server
	if err := http.Serve(socket, mux); err != nil {
		log.Fatal("Could not start HTTP server.")
	}
}

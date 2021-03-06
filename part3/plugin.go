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

// Generic struct for HTTP json requests from the Docker daemon
// The `json:"Name"` annotation tell the JSON decoder which fields
// correspond to the struct variables.
type request struct {
	Name string            `json:"Name"`
	Opts map[string]string `json:"Opts"`
}

// Generic struct for plugin response. The `omitempty` annotation prevents
// the struct variable to be encoded in the output if it is empty.
type pluginResponse struct {
	Implements []string `json:"Implements,omitempty"`
	Err        string   `json:"Err,omitempty"`
}

// Convenience type to easily create http handlers. We wrap every route function
// in this type to handle the boilerplate of error handling and JSON decoding/encoding.
type requestHandler func(request) interface{}

// Every handler should have a ServeHTTP function to be used by the HTTP mux
// This is a generic implementation that is used by every route
func (rh requestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req request

	// Decode the JSON request into a struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, pluginResponse{Err: "Malformed request: " + err.Error()})
		return
	}

	log.Printf("Received request for url: %s, with body: %+v", r.URL, req)

	// Call the appropriate handler and send response back
	sendResponse(w, rh(req))
}

// Send a json response with the correct headers and JSON structure
func sendResponse(w http.ResponseWriter, res interface{}) {
	// Docker expects to find a header with the correct mime-type and version
	w.Header().Set("Content-Type", versionMimetype)

	// Encode the JSON repsonse body
	json.NewEncoder(w).Encode(res)
}

// Activate the plugin with the Docker daemon
func activate(req request) interface{} {
	return pluginResponse{Implements: []string{"VolumeDriver"}}
}

// Create the volume
func create(req request) interface{} {
	// return empty struct
	return struct{}{}
}

// Remove the volume
func remove(req request) interface{} {
	// return empty struct
	return struct{}{}
}

// Mount the volume on the filesystem
func mount(req request) interface{} {
	// return empty struct
	return struct{}{}
}

// Unmount the volume
func unmount(req request) interface{} {
	// return empty struct
	return struct{}{}
}

// Return the path to the volume
func path(req request) interface{} {
	// return empty struct
	return struct{}{}
}

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

	// Add plugin and driver routes. Wrap every route function in a
	// requestHandler type so that it can be used by the HTTP muxer.
	mux.Handle("/Plugin.Activate", requestHandler(activate))
	mux.Handle("/VolumeDriver.Create", requestHandler(create))
	mux.Handle("/VolumeDriver.Remove", requestHandler(remove))
	mux.Handle("/VolumeDriver.Mount", requestHandler(mount))
	mux.Handle("/VolumeDriver.Unmount", requestHandler(unmount))
	mux.Handle("/VolumeDriver.Path", requestHandler(path))

	log.Print("Starting Dutch Docker Day volume plugin...")

	// Start the HTTP server
	if err := http.Serve(socket, mux); err != nil {
		log.Fatal("Could not start HTTP server.")
	}
}

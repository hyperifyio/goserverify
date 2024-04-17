// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	v8 "rogchap.com/v8go"
)

func main() {

	staticDir := "./frontend/build"
	serverPort := "8080"

	setupHttpServer(
		serverPort,
		staticDir,
	)

}

func setupHttpServer(
	serverPort,
	staticDir string,
) {

	router := mux.NewRouter()

	// Middleware to handle static files if they exist
	fileServer := http.FileServer(http.Dir(staticDir))
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			// If the path is not valid, directly go to dynamic handling
			dynamicHandler(w, r)
			return
		}
		// Check if the file exists and or if it is a directory (which we also treat as not found)
		filePath := filepath.Join(staticDir, path)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r) // Serve static files
		} else {
			dynamicHandler(w, r) // Handle dynamically
		}
	})

	log.Printf("Server is running at http://localhost:%s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, router))

}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {

	// JavaScript code with embedded HTML/CSS
	jsCode := `
        var html = '<html><head><title>Test Page</title><style>body { background-color: #f0f0f0; }</style></head><body><h1>Hello from Go and V8!</h1></body></html>';
        html;
    `

	// Initialize V8 context
	ctx := v8.NewContext()
	defer ctx.Close() // Ensure the context is properly disposed of

	// Execute JavaScript
	val, err := ctx.RunScript(jsCode, "render.js")
	if err != nil {
		log.Printf("Failed to execute script: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		// Send the rendered HTML as the response
		fmt.Fprint(w, val)
	}

}

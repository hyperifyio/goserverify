// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net/http"
	v8 "rogchap.com/v8go"
)

func main() {

	staticDir := "./frontend/build"
	staticPath := "/static"

	serverPort := "8080"

	// Handle dynamic requests
	http.HandleFunc("/", handler)

	// Serve static files from the "./static" directory
	fs := http.FileServer(http.Dir(staticDir))

	// Use the "/static/" prefix to access static content
	http.Handle(staticPath+"/", http.StripPrefix(staticPath+"/", fs))

	log.Printf("Server is running at http://localhost:%s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

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

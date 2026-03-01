package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Serve static files from the assets directory
	fs := http.FileServer(http.Dir("./cmd/web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Serve the main index page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./cmd/web/index.html")
	})

	fmt.Printf("Old Man Supper Club server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

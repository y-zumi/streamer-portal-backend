package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// main run http server
func main() {
	log.Print("starting server...")

	http.HandleFunc("/live_statuses", listLiveStatusesHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

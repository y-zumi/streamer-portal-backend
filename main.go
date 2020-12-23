package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ListLiveStatusesRequest is request to get live statuses
type ListLiveStatusesRequest struct {
	StreamerID string `json:"streamer_id"`
}

// ListLiveStatusesResponse is response to get live statuses
type ListLiveStatusesResponse struct {
	LiveStatuses []LiveStatus `json:"live_statuses"`
}

// LiveStatus represents streamer's live status on the stream platforms
type LiveStatus struct {
	PlatformType string `json:"platform_type"`
	IsLive       bool   `json:"is_live"`
}

// handler handle live status endpoint
func handler(w http.ResponseWriter, r *http.Request) {
	var req ListLiveStatusesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decode error: %v", err)
		return
	}
	fmt.Printf("req: %v", req)

	resp := ListLiveStatusesResponse{
		LiveStatuses: []LiveStatus{
			{
				PlatformType: "youtube",
				IsLive:       true,
			},
			{
				PlatformType: "twitch",
				IsLive:       false,
			},
			{
				PlatformType: "niconico",
				IsLive:       false,
			},
		},
	}

	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("failed to marshal resp: %v", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if _, err := w.Write(b); err != nil {
		fmt.Printf("failed to write resp: %v", err)
		return
	}
}

// main run http server
func main() {
	log.Print("starting server...")

	http.HandleFunc("/live_statuses", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

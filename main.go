package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GetLiveStatusesRequest struct {
	StreamerID string `json:"streamer_id"`
}

type GetLiveStatusesResponse struct {
	LiveStatuses []LiveStatus `json:"live_statuses"`
}

type LiveStatus struct {
	PlatformType string `json:"platform_type"`
	IsLive       bool   `json:"is_live"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var req GetLiveStatusesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decode error: %v", err)
		return
	}
	fmt.Printf("req: %v", req)

	resp := GetLiveStatusesResponse{
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

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/live_statuses", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

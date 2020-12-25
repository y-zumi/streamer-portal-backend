package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

// env has environment variables
type env struct {
	YoutubeAPIKey string `envconfig:"YOUTUBE_API_KEY"`
}

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

// listLiveStatusesHandler handle live status endpoint
func listLiveStatusesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var e env
	envconfig.Process("", &e)

	var req ListLiveStatusesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decode error: %v", err)
		return
	}

	client := YoutubeAPIClient{apiKey: e.YoutubeAPIKey}
	youtube, err := client.GetLiveStatus(ctx, req.StreamerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp := ListLiveStatusesResponse{
		LiveStatuses: []LiveStatus{
			{
				PlatformType: "youtube",
				IsLive:       youtube.IsLive,
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

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
	YoutubeAPIKey   string `envconfig:"YOUTUBE_API_KEY"`
	TwitchAuthToken string `envconfig:"TWITCH_AUTH_TOKEN"`
	TwitchClientID  string `envconfig:"TWITCH_CLIENT_ID"`
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
	if err := envconfig.Process("", &e); err != nil {
		fmt.Println(err)
		return
	}

	var req ListLiveStatusesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decode error: %v", err)
		return
	}

	client := YoutubeAPIClient{apiKey: e.YoutubeAPIKey}
	youtube, err := client.GetLive(ctx, "UCx1nAvtVDIsaGmCMSe8ofsQ")
	if err != nil {
		fmt.Println(err)
		return
	}

	twitchClient := TwitchClient{
		AuthToken: e.TwitchAuthToken,
		ClientID:  e.TwitchClientID,
	}
	twitch, err := twitchClient.GetLive(ctx, "545050196")
	if err != nil {
		fmt.Println(err)
		return
	}

	niconicoClient := NiconicoClient{}
	niconico, err := niconicoClient.GetLive(ctx, "2598430")
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
				IsLive:       twitch.IsLive,
			},
			{
				PlatformType: "niconico",
				IsLive:       niconico.IsLive,
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

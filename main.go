package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

// YoutubeAPIClient calls Youtube Data API
type YoutubeAPIClient struct {
	apiKey string
}

// searchResponse represents youtube api /search endpoint's response
type searchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

// Live represents youtube live information
type Live struct {
	IsLive  bool
	VideoID string
}

const (
	// baseUrl represents Youtube Data API endpoint
	baseUrl = "https://www.googleapis.com/youtube/v3"
)

// GetLiveStatus get streamer's live status by channel ID in Youtube
func (c *YoutubeAPIClient) GetLiveStatus(ctx context.Context, channelID string) (*Live, error) {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 10,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/search?part=snippet&eventType=live&type=video&fields=items(snippet/liveBroadcastContent,id/videoId)&channelId=%s&key=%s", baseUrl, channelID, c.apiKey),
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}
	defer resp.Body.Close()

	s := new(searchResponse)
	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	if len(s.Items) == 0 {
		return &Live{
			IsLive: false,
		}, nil
	}

	return &Live{
		IsLive:  s.Items[0].Snippet.LiveBroadcastContent == "live",
		VideoID: s.Items[0].ID.VideoID,
	}, nil
}

// handler handle live status endpoint
func handler(w http.ResponseWriter, r *http.Request) {
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

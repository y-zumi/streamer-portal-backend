package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

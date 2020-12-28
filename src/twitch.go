package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TwitchClient struct {
	AuthToken string
	ClientID  string
}

type twitchStreamsResponse struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		UserName     string    `json:"user_name"`
		GameID       string    `json:"game_id"`
		GameName     string    `json:"game_name"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailURL string    `json:"thumbnail_url"`
		TagIDs       []string  `json:"tag_ids"`
	} `json:"data"`
}

const (
	twitchBaseUrl = "https://api.twitch.tv/helix"
)

func (c *TwitchClient) GetLive(ctx context.Context, userID string) (*Live, error) {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/streams?user_id=%s", twitchBaseUrl, userID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", c.AuthToken)},
		"Client-ID":     {c.ClientID},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var videos twitchStreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&videos); err != nil {
		return nil, err
	}
	if len(videos.Data) == 0 {
		return &Live{
			IsLive:  false,
			VideoID: "",
		}, nil
	}

	return &Live{
		IsLive:  videos.Data[0].Type == "live",
		VideoID: videos.Data[0].ID,
	}, nil
}

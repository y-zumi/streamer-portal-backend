package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type TwitchClient struct {
	AuthToken string
	ClientID  string
}

const (
	twitchBaseUrl = "https://api.twitch.tv/helix/videos"
)

func (t *TwitchClient) GetLiveStatus(ctx context.Context, userID string) (*Live, error) {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?user_id=%s", twitchBaseUrl, userID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", t.AuthToken)},
		"Client-ID":     {t.ClientID},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)

	return nil, nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type NiconicoClient struct{}

type niconicoLiveResponse struct {
	Data []struct {
		StartTime      string `json:"startTime"`
		CategoryTags   string `json:"categoryTags"`
		ContentId      string `json:"contentId"`
		CommentCounter int64  `json:"commentCounter"`
		ChannelId      int64  `json:"channelId"`
		LiveStatus     string `json:"liveStatus"`
		Description    string `json:"description"`
		Tags           string `json:"tags"`
		UserId         int64  `json:"userId"`
		Title          string `json:"title"`
	} `json:"data"`
}

const (
	niconicoBaseUrl = "https://api.search.nicovideo.jp/api/v2"
)

func (c NiconicoClient) GetLive(ctx context.Context, channelID string) (*Live, error) {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/live/contents/search?targets=title,description,tags,tagsExact,categoryTags&_sort=-startTime&fields=title,description,channelId,commentCounter,userId,categoryTags,contentId,tags,liveStatus,startTime&q=%s&filters[liveStatus][0]=onair&filters[channelId][0]=%s", niconicoBaseUrl, url.QueryEscape("一般(その他) OR ゲーム"), channelID),
		nil,
	)
	if err != nil {
		fmt.Println("req")
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp")
		return nil, err
	}

	var live niconicoLiveResponse
	if err := json.NewDecoder(resp.Body).Decode(&live); err != nil {
		fmt.Println("decode")
		return nil, err
	}
	if len(live.Data) == 0 {
		return &Live{
			IsLive:  false,
			VideoID: "",
		}, nil
	}

	return &Live{
		IsLive:  live.Data[0].LiveStatus == "onair",
		VideoID: live.Data[0].ContentId,
	}, nil
}

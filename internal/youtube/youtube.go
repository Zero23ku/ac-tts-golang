package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var API_KEY = ""
var VIDEO_ID = ""
var livestreamChatId = ""

type Snippet struct {
	ChannelId            string `json:"channelId"`
	LiveBroadcastContent string `json:"liveBroadcastContent"`
}

type LiveStreamingDetails struct {
	ActiveLiveChatId string `json:"activeLiveChatId"`
}

type Item struct {
	Kind                 string               `json:"kind"`
	Etag                 string               `json:"etag"`
	Id                   string               `json:"id"`
	Snippet              Snippet              `json:"snippet"`
	LiveStreamingDetails LiveStreamingDetails `json:"liveStreamingDetails"`
}

type YTChannelInfo struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []Item `json:"items"`
}

type TextMessageDetails struct {
	MessageText string `json:"messageText"`
}

type SnippetChat struct {
	Type               string             `json:"type"`
	TextMessageDetails TextMessageDetails `json:"textMessageDetails"`
}

type ItemChat struct {
	Kind    string      `json:"kind"`
	Etag    string      `json:"etag"`
	Id      string      `json:"id"`
	Snippet SnippetChat `json:"snippet"`
}

type LivechatResponse struct {
	Kind                  string     `json:"kind"`
	Etag                  string     `json:"etag"`
	NextpageToken         string     `json:"nextPageToken"`
	Items                 []ItemChat `json:"items"`
	PollingIntervalMillis int        `json:"pollingIntervalMillis"`
}

const liveStreamingDetailsEndpoint = "https://www.googleapis.com/youtube/v3/videos"

const liveStreamingGetChatMessages = "https://www.googleapis.com/youtube/v3/liveChat/messages"

func GetYTChannelInfo() {

	client := &http.Client{}
	url := liveStreamingDetailsEndpoint + "?part=liveStreamingDetails,snippet&id=" + VIDEO_ID + "&key=" + API_KEY
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error sending request", err)
	}
	defer resp.Body.Close()

	var response YTChannelInfo

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}

	//TODO: VALIDAR QUE Items no esté vacio
	livestreamChatId = response.Items[0].LiveStreamingDetails.ActiveLiveChatId

	go func() {

		client := &http.Client{}
		pageToken := ""
		chatUrl := liveStreamingGetChatMessages + "?liveChatId=" + livestreamChatId + "&part=snippet,authorDetails&maxResults=1000&key=" + API_KEY

		for {

			if pageToken != "" {
				chatUrl = chatUrl + "&pageToken=" + pageToken
			}

			req, err := http.NewRequest("GET", chatUrl, nil)

			if err != nil {
				log.Fatal(err)
			}

			resp, err := client.Do(req)

			if err != nil {
				log.Fatal("Error sending request", err)
			}
			defer resp.Body.Close()

			var response LivechatResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				log.Fatal(err)
			}

			for i := 0; i < len(response.Items); i++ {
				//TODO: Hacer sonar los mensajes
				fmt.Println(response.Items[i].Snippet.TextMessageDetails.MessageText)
			}

			pageToken = response.NextpageToken
			interval := response.PollingIntervalMillis
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}

	}()

}

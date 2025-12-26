package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"ac-tts/internal/common"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
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

var ytWindowIsOpen = false

const liveStreamingDetailsEndpoint = "https://www.googleapis.com/youtube/v3/videos"

const liveStreamingGetChatMessages = "https://www.googleapis.com/youtube/v3/liveChat/messages"

var YoutubeWindow fyne.Window
var ConnectYTButton *widget.Button
var AppReference *fyne.App
var CTX context.Context

func GetYTChannelInfo(ctx context.Context) {

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

	if len(response.Items) > 0 {
		livestreamChatId = response.Items[0].LiveStreamingDetails.ActiveLiveChatId
	}

	go func() {

		client := &http.Client{}
		pageToken := ""
		chatUrl := liveStreamingGetChatMessages + "?liveChatId=" + livestreamChatId + "&part=snippet,authorDetails&maxResults=1000&key=" + API_KEY
		for {

			select {
			case <-ctx.Done():
				return
			default:
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
					ytMsg := response.Items[i].Snippet.TextMessageDetails.MessageText
					if common.IsTTSCommandActive() && strings.HasPrefix(ytMsg, common.GetTTSCommand()) {
						reproductor.Reproduce(ytMsg, "")
					} else if !common.IsTTSCommandActive() {
						reproductor.Reproduce(ytMsg, "")
					}

					time.Sleep(time.Duration(1200) * time.Millisecond)
				}

				pageToken = response.NextpageToken
				interval := response.PollingIntervalMillis
				time.Sleep(time.Duration(interval) * time.Millisecond)
			}

		}

	}()

}

func initYoutubeWindow(app fyne.App) {
	YoutubeWindow = app.NewWindow("Youtube Integration (Alpha)")
	YoutubeWindow.SetOnClosed(func() {
		ytWindowIsOpen = false
	})

	ytApiKeyInput := widget.NewEntry()
	ytApiKeyInput.SetPlaceHolder("Enter your Youtube's API Key here")
	ytApiKeyInput.Resize(fyne.NewSize(100, ytApiKeyInput.MinSize().Height))

	ytVideoInput := widget.NewEntry()
	ytVideoInput.SetPlaceHolder("Enter Livestream's url: https://www.youtube.com/watch?v=your-id")
	ytVideoInput.Resize(fyne.NewSize(100, ytVideoInput.MinSize().Height))

	ytApiKeySubmit := widget.NewButton("Connect", func() {
		API_KEY = ytApiKeyInput.Text

		ytID, err := getYTID(ytVideoInput.Text)
		if err != nil && ytID == "" {
			logging.CreateLog("Youtube - ", err)
		}
		VIDEO_ID = ytID
		GetYTChannelInfo(CTX)
		YoutubeWindow.Close()
	})

	form := widget.NewForm(
		widget.NewFormItem("Youtube's API Key", ytApiKeyInput),
	)

	formVide := widget.NewForm(
		widget.NewFormItem("Youtube livestream URL", ytVideoInput),
	)

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		ytApiKeySubmit,
	)

	YoutubeWindow.SetContent(container.NewVBox(form, formVide, centeredButton))
	YoutubeWindow.Resize(fyne.NewSize(400, 100))
}

func InitConnectYTButton() {
	ConnectYTButton = widget.NewButton("Connect to Youtube (Alpha)", func() {
		initYoutubeWindow(*AppReference)
		if !ytWindowIsOpen {
			YoutubeWindow.Show()
			ytWindowIsOpen = true
		}

	})
}

func getYTID(ytUrl string) (string, error) {
	parsedUrl, err := url.Parse(ytUrl)
	if err != nil {
		return ytUrl, err
	}

	qParams := parsedUrl.Query()
	id := qParams.Get("v")
	if id == "" {
		return "", fmt.Errorf("No video ID")
	}

	return id, nil
}

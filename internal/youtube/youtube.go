package youtube

import (
	"context"
	"encoding/json"
	"errors"
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
	YtChat "github.com/johanvandegriff/youtube-live-chat-downloader/v2"

	"ac-tts/internal/common"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/whitelist"
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
	Kind          string        `json:"kind"`
	Etag          string        `json:"etag"`
	Id            string        `json:"id"`
	Snippet       SnippetChat   `json:"snippet"`
	AuthorDetails AuthorDetails `json:"authorDetails"`
}

type LivechatResponse struct {
	Kind                  string     `json:"kind"`
	Etag                  string     `json:"etag"`
	NextpageToken         string     `json:"nextPageToken"`
	Items                 []ItemChat `json:"items"`
	PollingIntervalMillis int        `json:"pollingIntervalMillis"`
}

type AuthorDetails struct {
	DisplayName string `json:"displayName"`
}

var ytWindowIsOpen = false

const liveStreamingDetailsEndpoint = "https://www.googleapis.com/youtube/v3/videos"

const liveStreamingGetChatMessages = "https://www.googleapis.com/youtube/v3/liveChat/messages"

var YoutubeWindow fyne.Window
var ConnectYTButton *widget.Button
var AppReference *fyne.App
var CTX context.Context
var YoutubeErrorWindow fyne.Window

// No se usa, pero lo dejo en caso de que la librería que se usa actualmente para el chat dejase de funcionar
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
					initErrorWindow(*AppReference, err.Error())
					YoutubeErrorWindow.Show()
					return
				}

				resp, err := client.Do(req)

				if err != nil {
					initErrorWindow(*AppReference, err.Error())
					YoutubeErrorWindow.Show()
					return
				}
				defer resp.Body.Close()
				var response LivechatResponse

				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					initErrorWindow(*AppReference, err.Error())
					YoutubeErrorWindow.Show()
					return
				}
				for i := 0; i < len(response.Items); i++ {

					ytMsg := response.Items[i].Snippet.TextMessageDetails.MessageText
					chatter := response.Items[i].AuthorDetails.DisplayName
					if common.IsTTSCommandActive() && strings.HasPrefix(ytMsg, common.GetTTSCommand()) {
						if common.IsPitchRandom {
							reproductor.Reproduce(ytMsg, chatter, common.GetRandomPitch(), common.IsPitchRandom)
						} else {
							reproductor.Reproduce(ytMsg, chatter, common.Pitch, common.IsPitchRandom)
						}
					} else if !common.IsTTSCommandActive() {
						if common.IsPitchRandom {
							reproductor.Reproduce(ytMsg, chatter, common.GetRandomPitch(), common.IsPitchRandom)
						} else {
							reproductor.Reproduce(ytMsg, chatter, common.Pitch, common.IsPitchRandom)
						}
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
	/*
		ytApiKeyInput := widget.NewEntry()
		ytApiKeyInput.SetPlaceHolder("Enter your Youtube's API Key here")
		ytApiKeyInput.Resize(fyne.NewSize(100, ytApiKeyInput.MinSize().Height))
	*/

	ytVideoInput := widget.NewEntry()
	ytVideoInput.SetPlaceHolder("Enter Livestream's url: https://www.youtube.com/watch?v=your-id")
	ytVideoInput.Resize(fyne.NewSize(100, ytVideoInput.MinSize().Height))

	ytApiKeySubmit := widget.NewButton("Connect", func() {
		//API_KEY = ytApiKeyInput.Text
		param, err := validateYTUrl(ytVideoInput.Text)
		if err != nil {
			initErrorWindow(app, err.Error())
			return
		}

		ytUrl := "https://www.youtube.com/watch?v=" + param

		ytID, err := getYTID(ytUrl)
		if err != nil && ytID == "" {
			logging.CreateLog("Youtube - ", err)
		}
		VIDEO_ID = ytID
		//GetYTChannelInfo(CTX)
		connectToChat()
		YoutubeWindow.Close()
	})

	/*form := widget.NewForm(
		widget.NewFormItem("Youtube's API Key", ytApiKeyInput),
	)*/

	formVide := widget.NewForm(
		widget.NewFormItem("Youtube livestream URL", ytVideoInput),
	)

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		ytApiKeySubmit,
	)

	YoutubeWindow.SetContent(container.NewVBox( /*form,*/ formVide, centeredButton))
	YoutubeWindow.Resize(fyne.NewSize(400, 100))
}

func validateYTUrl(rawUrl string) (string, error) {
	url, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return "", err
	}
	hostname := strings.Split(url.Host, ".")[1]
	if hostname != "youtube" {
		return "", errors.New("URl must be from Youtube")
	}

	param := url.Query().Get("v")
	if param == "" {
		path := strings.Split(url.Path, "/")
		param = path[len(path)-1]
	}
	return param, nil

}

func InitConnectYTButton() {
	ConnectYTButton = widget.NewButton("Connect to Youtube", func() {
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

func connectToChat() {

	continuation, cfg, error := YtChat.ParseInitialData("https://www.youtube.com/watch?v=" + VIDEO_ID)
	if error != nil {
		logging.CreateLog("Youtube getting chat- ", error)
		initErrorWindow(*AppReference, "Error connecting to livestream, please try again.")
		YoutubeErrorWindow.Show()
		//log.Fatal(error)
	}
	go func() {
		for {
			select {
			case <-CTX.Done():
				return
			default:
				chat, newContinuation, error := YtChat.FetchContinuationChat(continuation, cfg)
				if error == YtChat.ErrLiveStreamOver {
					logging.CreateLog("Youtube livestream - ", error)
					initErrorWindow(*AppReference, "Livestream is over")
					YoutubeErrorWindow.Show()
					return

				}
				if error != nil {
					initErrorWindow(*AppReference, error.Error())
					YoutubeErrorWindow.Show()
					return
				}
				continuation = newContinuation

				for _, msg := range chat {
					/*fmt.Print(msg.Timestamp, " | ")*/
					//fmt.Println(msg.AuthorName, ": ", msg.Message)
					author := msg.AuthorName[1:]
					if common.IsTTSCommandActive() && strings.HasPrefix(msg.Message, common.GetTTSCommand()) {
						if common.IsPitchRandom {
							if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(author) {
								reproductor.Reproduce(msg.Message, author, common.GetRandomPitch(), common.IsPitchRandom)
							} else if !whitelist.IsWhitelistActive {
								reproductor.Reproduce(msg.Message, author, common.GetRandomPitch(), common.IsPitchRandom)
							}

						} else {
							if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(author) {
								reproductor.Reproduce(msg.Message, author, common.Pitch, common.IsPitchRandom)
							} else if !whitelist.IsWhitelistActive {
								reproductor.Reproduce(msg.Message, author, common.Pitch, common.IsPitchRandom)
							}
						}
					} else if !common.IsTTSCommandActive() {
						if common.IsPitchRandom {
							if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(author) {
								reproductor.Reproduce(msg.Message, author, common.GetRandomPitch(), common.IsPitchRandom)
							} else if !whitelist.IsWhitelistActive {
								reproductor.Reproduce(msg.Message, author, common.GetRandomPitch(), common.IsPitchRandom)
							}
						} else {
							if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(author) {
								reproductor.Reproduce(msg.Message, author, common.Pitch, common.IsPitchRandom)
							} else if !whitelist.IsWhitelistActive {
								reproductor.Reproduce(msg.Message, author, common.Pitch, common.IsPitchRandom)
							}
						}
					}

					time.Sleep(1000 * time.Millisecond)
				}
			}

		}
	}()
}

func initErrorWindow(app fyne.App, msg string) {
	YoutubeErrorWindow = app.NewWindow("Error!")
	YoutubeErrorWindow.SetContent(widget.NewLabel(msg))
}

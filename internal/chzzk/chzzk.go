package chzzk

import (
	"context"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dhkimxx/GoChzzkChatCrawler/crawler"

	"ac-tts/internal/common"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/whitelist"
)

var chzzkWindowIsOpen = false
var CHZZKWindow fyne.Window
var CHZZKErrorWindow fyne.Window
var CHZZKAlertWindow fyne.Window
var ConnectChzzkButton *widget.Button
var AppReference *fyne.App
var CTX context.Context
var chzzkLiveId = ""

func connectToChzzk(liveid string, ctx context.Context) {

	// Create a new crawler client with callback handler
	crawlerClient := crawler.NewCrawlerClient(liveid, 1, func(msg crawler.ChzzkChatMessage) {

		if common.IsTTSCommandActive() && strings.HasPrefix(msg.Content, common.GetTTSCommand()) {
			if common.IsPitchRandom {
				if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(msg.Nickname) {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.GetRandomPitch(), common.IsPitchRandom)
				} else if !whitelist.IsWhitelistActive {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.GetRandomPitch(), common.IsPitchRandom)
				}

			} else {

				if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(msg.Nickname) {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.Pitch, common.IsPitchRandom)
				} else if !whitelist.IsWhitelistActive {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.Pitch, common.IsPitchRandom)
				}
			}
		} else if !common.IsTTSCommandActive() {
			if common.IsPitchRandom {
				if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(msg.Nickname) {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.GetRandomPitch(), common.IsPitchRandom)
				} else if !whitelist.IsWhitelistActive {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.GetRandomPitch(), common.IsPitchRandom)
				}
			} else {
				if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(msg.Nickname) {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.Pitch, common.IsPitchRandom)
				} else if !whitelist.IsWhitelistActive {
					reproductor.Reproduce(msg.Content, msg.Nickname, common.Pitch, common.IsPitchRandom)
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	})

	go func() {
		// Start crawling
		err := crawlerClient.Run()
		if err != nil {
			CHZZKErrorWindow.Show()
			return
		}
		for {
			select {
			case <-ctx.Done():
				return

			}
		}
	}()

}

func initChzzkWindow(app fyne.App) {
	CHZZKWindow = app.NewWindow("chzzk integration")
	initChzzkErrorWindow(app)
	CHZZKWindow.SetOnClosed(func() {
		chzzkWindowIsOpen = false
	})

	chzzkChannelNameInput := widget.NewEntry()
	chzzkChannelNameInput.SetPlaceHolder("Enter chzzk live id")
	chzzkChannelNameInput.Resize(fyne.NewSize(100, chzzkChannelNameInput.MinSize().Height))

	chzzkSubmit := widget.NewButton("Connect", func() {
		chzzkLiveId = chzzkChannelNameInput.Text
		if chzzkLiveId == "" {
			initChzzkAlertWindow(*AppReference, "You must enter a livestream id")
		} else {
			connectToChzzk(chzzkLiveId, CTX)
			CHZZKWindow.Close()
		}
	})

	form := widget.NewForm(
		widget.NewFormItem("chzzk live's id", chzzkChannelNameInput),
	)

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		chzzkSubmit,
	)

	CHZZKWindow.SetContent(
		container.NewVBox(form, centeredButton),
	)
	CHZZKWindow.Resize(fyne.NewSize(400, 400))
}

func InitConnectChzzkButton() {
	ConnectChzzkButton = widget.NewButton("Connect to chzzk", func() {
		initChzzkWindow(*AppReference)
		if !chzzkWindowIsOpen {
			CHZZKWindow.Show()
			chzzkWindowIsOpen = true
		}
	})
}

func initChzzkErrorWindow(app fyne.App) {
	CHZZKErrorWindow = app.NewWindow("Error!")
	CHZZKErrorWindow.SetContent(widget.NewLabel("An error ocurrier while connecting to chzzk, please try again!"))
}

func initChzzkAlertWindow(app fyne.App, msg string) {
	CHZZKErrorWindow = app.NewWindow("Error!")
	CHZZKErrorWindow.SetContent(widget.NewLabel(msg))
}

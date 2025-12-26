package tiktok

import (
	"context"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/steampoweredtaco/gotiktoklive"

	"ac-tts/internal/common"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
)

var tiktokWindowIsOpen = false

var TiktokWindow fyne.Window
var TiktokErrorWindow fyne.Window
var ConnectTiktokButton *widget.Button
var AppReference *fyne.App
var CTX context.Context
var tiktokChannelName = ""

func connectToTikTokChat(username string, ctx context.Context) {
	//fmt.Println("Iniciando tiktok...")
	tiktok, err := gotiktoklive.NewTikTok()
	if err != nil {
		logging.CreateLog("Error initializing Tiktok..", err)
		log.Fatal(err)
	}
	//fmt.Println("Capturando chat: ", username)
	live, err := tiktok.TrackUser(username)

	if err != nil {
		TiktokErrorWindow.Show()
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-live.Events:
				switch e := event.(type) {
				case gotiktoklive.ChatEvent:
					tiktokMsg := gotiktoklive.ChatEvent(e).Comment
					if common.IsTTSCommandActive() && strings.HasPrefix(tiktokMsg, common.GetTTSCommand()) {
						reproductor.Reproduce(tiktokMsg, "")
					} else if !common.IsTTSCommandActive() {
						reproductor.Reproduce(tiktokMsg, "")
					}

					time.Sleep(time.Duration(1200) * time.Millisecond)
				}
			}
		}
	}()
}

func initTiktokWindow(app fyne.App) {
	TiktokWindow = app.NewWindow("Tiktok integration")
	initTiktokErrorWindow(app)
	TiktokWindow.SetOnClosed(func() {
		tiktokWindowIsOpen = false
	})

	tiktokChannelNameInput := widget.NewEntry()
	tiktokChannelNameInput.SetPlaceHolder("Enter Tiktok channel's name")
	tiktokChannelNameInput.Resize(fyne.NewSize(100, tiktokChannelNameInput.MinSize().Height))

	tiktokSubmit := widget.NewButton("Connect", func() {
		tiktokChannelName = tiktokChannelNameInput.Text
		connectToTikTokChat(tiktokChannelName, CTX)
		TiktokWindow.Close()
	})

	form := widget.NewForm(
		widget.NewFormItem("Tiktok Channel's name", tiktokChannelNameInput),
	)

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		tiktokSubmit,
	)

	TiktokWindow.SetContent(container.NewVBox(form, centeredButton))
	TiktokWindow.Resize(fyne.NewSize(400, 100))

}

func InitConnectTiktokButton() {
	ConnectTiktokButton = widget.NewButton("Connect to Tiktok", func() {
		initTiktokWindow(*AppReference)
		if !tiktokWindowIsOpen {
			TiktokWindow.Show()
			tiktokWindowIsOpen = true
		}
	})
}

func initTiktokErrorWindow(app fyne.App) {
	TiktokErrorWindow = app.NewWindow("Error!")
	TiktokErrorWindow.SetContent(widget.NewLabel("An error ocurried while connecting to Tiktok, please try again."))
}

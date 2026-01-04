package common

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"ac-tts/internal/assets"
	"ac-tts/internal/config"
)

var TwitchConnectButton *widget.Button
var PitchSlider *widget.Slider
var PitchRow *fyne.Container
var KofiButton *widget.Button
var TestPitchButton *widget.Button
var UpdateButton *widget.Button
var ActivateCommand *widget.Check
var InputCommand *widget.Entry
var RandomPitch *widget.Check
var DocsButton *widget.Button
var AppReference *fyne.App

// Internal
var leftSpacer *canvas.Rectangle
var left *fyne.Container
var pitchLabel *canvas.Text
var kofiUrl *url.URL
var docsUrl *url.URL
var githubUrl *url.URL
var isCommandActive = false
var ttsCommand = "!tts" //valor default

// Twitch conf

var TwitchConfWindow fyne.Window
var IsRedeemOptionActiva = false
var TwitchRedeemName *widget.Entry
var ConnectToTwitch *widget.Button
var IsTwitchConnected = false
var TwitchActiveRedeemOption *widget.Check
var twitchConfiWindowsIsOpen = false
var TwitchErrorWindow fyne.Window

func initTwitchConfWindow(app fyne.App, onClick func()) {
	TwitchConfWindow = app.NewWindow("Twitch configuration")

	TwitchConfWindow.SetOnClosed(func() {
		twitchConfiWindowsIsOpen = false
	})

	TwitchRedeemName = widget.NewEntry()
	TwitchRedeemName.SetPlaceHolder("Enter Redeem's name")
	TwitchRedeemName.Disable()
	TwitchRedeemName.Resize(fyne.NewSize(100, TwitchRedeemName.MinSize().Height))

	form := widget.NewForm(
		widget.NewFormItem("Twitch Redeem's name", TwitchRedeemName),
	)

	radio := widget.NewRadioGroup([]string{"Read all messages in chat", "Use channel points"}, func(value string) {
		if value == "Read all messages in chat" {
			IsRedeemOptionActiva = false
			TwitchRedeemName.Disable()
		} else {
			IsRedeemOptionActiva = true
			TwitchRedeemName.Enable()
		}
		form.Refresh()
	})

	radio.SetSelected("Read all messages in chat")

	ConnectToTwitch = widget.NewButton("Connect to Twitch", func() {
		redeemText := ""
		if IsRedeemOptionActiva {
			redeemText = TwitchRedeemName.Text
		}
		config.SaveConfig(!IsRedeemOptionActiva, redeemText)
		if IsRedeemOptionActiva && TwitchRedeemName.Text == "" {
			initTwitchErrorWindow(app, "You must enter a Redeem's name")
			TwitchErrorWindow.Show()
			return
		}
		onClick()
		TwitchConfWindow.Close()
	})

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		ConnectToTwitch,
	)

	TwitchConfWindow.SetContent(
		container.NewVBox(
			radio,
			form,
			centeredButton,
		),
	)

	TwitchConfWindow.Resize(fyne.NewSize(400, 100))
	configs, err := config.ReadConfig()
	if err != nil {
		initTwitchErrorWindow(app, "Error reading tts_config.json")
		TwitchErrorWindow.Show()
	} else {
		if configs.TwitchConfig.ReadAllChat {
			radio.SetSelected("Read all messages in chat")
		} else {
			radio.SetSelected("Use channel points")
			TwitchRedeemName.SetText(configs.TwitchConfig.RedeemName)
		}
	}

}

func InitLeftSpacer() {
	leftSpacer = canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(20, 0))
}

func InitTestPitchButton(onClick func()) {
	TestPitchButton = widget.NewButton("Test Voice", onClick)
}

func InitTwitchConnectButton(onClick func()) {

	TwitchConnectButton = widget.NewButton("Connect to Twitch", func() {
		initTwitchConfWindow(*AppReference, onClick)
		if !twitchConfiWindowsIsOpen {
			TwitchConfWindow.Show()
			twitchConfiWindowsIsOpen = true
		}
	})
}

func SetConnected() {
	fyne.Do(func() {
		TwitchConnectButton.SetText("Connected")
		TwitchConnectButton.Disable()
	})
}

func initPitchSlider(pitchData binding.Float) {
	PitchSlider = widget.NewSliderWithData(0.4, 2.0, pitchData)
	PitchSlider.Step = 0.1
	PitchSlider.Resize(fyne.NewSize(500, 200))
}

func initLeftPitchLabel() {
	pitchLabel = canvas.NewText("Voice Pitch", theme.Color(theme.ColorNameForeground))
	left = container.NewHBox(leftSpacer, pitchLabel)
	PitchRow = container.New(
		layout.NewBorderLayout(nil, nil, left, nil),
		left,
		PitchSlider,
	)
}

func InitPitchRow(pitchData binding.Float) {
	initPitchSlider(pitchData)
	initLeftPitchLabel()
}

func InitKofiButton() {
	kofiUrl = &url.URL{
		Scheme: "https",
		Host:   "ko-fi.com",
		Path:   "/I3I41O6OUD",
	}

	res := fyne.NewStaticResource("cup-border.png", assets.Cup)

	KofiButton = widget.NewButtonWithIcon("Support me!", res, func() {
		fyne.CurrentApp().OpenURL(kofiUrl)
	})
}

// https://github.com/Zero23ku/ac-tts-golang/blob/main/docs/docs.md
func InitDocsButton() {
	docsUrl = &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "Zero23ku/ac-tts-golang/blob/main/docs/docs.md",
	}

	DocsButton = widget.NewButton("How to use", func() {
		fyne.CurrentApp().OpenURL(docsUrl)
	})
}

func InitUpdateButton() {

	githubUrl = &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/Zero23ku/ac-tts-golang/releases",
	}

	UpdateButton = widget.NewButton("New Version Avaible", func() {
		fyne.CurrentApp().OpenURL(githubUrl)
	})

}

func InitCommandCheck() {
	ActivateCommand = widget.NewCheck("Use Command", func(value bool) {
		isCommandActive = value
	})
}

func InitCommandInput() {
	InputCommand = widget.NewEntry()
	InputCommand.Text = ttsCommand
}

func InitRandomPitch() {
	RandomPitch = widget.NewCheck("Random Pitch per user", func(value bool) {
		IsPitchRandom = value
	})
}

func GetTTSCommand() string {
	return ttsCommand
}

func IsTTSCommandActive() bool {
	return isCommandActive
}

func initTwitchErrorWindow(app fyne.App, msg string) {
	TwitchErrorWindow = app.NewWindow("!Error")
	TwitchErrorWindow.SetContent(widget.NewLabel(msg))
}

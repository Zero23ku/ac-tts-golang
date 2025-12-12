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
)

var ConnectButton *widget.Button
var PitchSlider *widget.Slider
var PitchRow *fyne.Container
var KofiButton *widget.Button
var TestPitchButton *widget.Button
var UpdateButton *widget.Button
var ActivateCommand *widget.Check
var InputCommand *widget.Entry

// Internal
var leftSpacer *canvas.Rectangle
var left *fyne.Container
var pitchLabel *canvas.Text
var kofiUrl *url.URL
var githubUrl *url.URL
var isCommandActive = false
var ttsCommand = "!tts" //valor default

func InitLeftSpacer() {
	leftSpacer = canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(20, 0))
}

func InitTestPitchButton(onClick func()) {
	TestPitchButton = widget.NewButton("Test Voice", onClick)
}

func InitConnectButton(onClick func()) {
	ConnectButton = widget.NewButton("Connect to Twitch", onClick)
}

func SetConnected() {
	fyne.Do(func() {
		ConnectButton.SetText("Connected")
		ConnectButton.Disable()
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

func GetTTSCommand() string {
	return ttsCommand
}

func IsTTSCommandActive() bool {
	return isCommandActive
}

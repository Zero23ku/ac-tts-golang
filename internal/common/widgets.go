package common

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var ConnectButton *widget.Button
var PitchSlider *widget.Slider
var PitchRow *fyne.Container
var KofiButton *widget.Button

// Internal
var leftSpacer *canvas.Rectangle
var left *fyne.Container
var pitchLabel *canvas.Text
var kofiUrl *url.URL

func InitConnectButton(onClick func()) {
	ConnectButton = widget.NewButton("Connect to Twitch", onClick)
}

func SetConnected() {
	ConnectButton.SetText("Connected")
	ConnectButton.Disable()
}

func initPitchSlider(pitchData binding.Float) {
	PitchSlider = widget.NewSliderWithData(0.2, 2.0, pitchData)
	PitchSlider.Step = 0.1
	PitchSlider.Resize(fyne.NewSize(500, 200))
}

func initLeftPitchLabel() {
	leftSpacer = canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(20, 0))
	pitchLabel = canvas.NewText("Voice Pitch", color.White)
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

func InitKofiButton(res fyne.Resource) {
	kofiUrl = &url.URL{
		Scheme: "https",
		Host:   "ko-fi.com",
		Path:   "/I3I41O6OUD",
	}
	KofiButton = widget.NewButtonWithIcon("Support me!", res, func() {
		fyne.CurrentApp().OpenURL(kofiUrl)
	})
}

package common

import (
	"fmt"
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
	"ac-tts/internal/youtube"
)

var ConnectButton *widget.Button
var PitchSlider *widget.Slider
var PitchRow *fyne.Container
var KofiButton *widget.Button
var TestPitchButton *widget.Button
var UpdateButton *widget.Button
var YoutubeWindow fyne.Window
var ConnectYTButton *widget.Button
var AppReference *fyne.App

// Internal
var leftSpacer *canvas.Rectangle
var left *fyne.Container
var pitchLabel *canvas.Text
var kofiUrl *url.URL
var githubUrl *url.URL
var ytWindowIsOpen = false

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

func InitConnectYTButton() {
	ConnectYTButton = widget.NewButton("Connect to Youtube", func() {
		initYoutubeWindow(*AppReference)
		if !ytWindowIsOpen {
			YoutubeWindow.Show()
			ytWindowIsOpen = true
		}

	})
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

func initYoutubeWindow(app fyne.App) {
	YoutubeWindow = app.NewWindow("Youtube Integration")
	YoutubeWindow.SetOnClosed(func() {
		ytWindowIsOpen = false
	})

	ytApiKeyInput := widget.NewEntry()
	ytApiKeyInput.SetPlaceHolder("Enter your Youtube's API Key here")
	ytApiKeyInput.Resize(fyne.NewSize(100, ytApiKeyInput.MinSize().Height))

	ytVideoInput := widget.NewEntry()
	ytVideoInput.SetPlaceHolder("Enter Livestream's url: https://www.youtube.com/watch?v=your-id")
	ytVideoInput.Resize(fyne.NewSize(100, ytVideoInput.MinSize().Height))

	ytApiKeySubmit := widget.NewButton("Submit Key", func() {
		//TODO: Guardarlos
		fmt.Println(ytApiKeyInput.Text)
		fmt.Println(ytVideoInput.Text)
		youtube.API_KEY = ytApiKeyInput.Text
		youtube.VIDEO_ID = ytVideoInput.Text
		youtube.GetYTChannelInfo()
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

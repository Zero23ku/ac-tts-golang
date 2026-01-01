package local

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"ac-tts/internal/common"
	"ac-tts/internal/reproductor"
)

var pitchSlider *widget.Slider
var readButton *widget.Button
var textArea *widget.Entry

var LocalContainer *fyne.Container
var LocalWindow fyne.Window
var LocalWindowIsOpen = false
var LocalButton *widget.Button
var AppReference *fyne.App

func initLocalWindow(app fyne.App) {
	pitchData := binding.BindFloat(&common.Pitch)
	pitchSlider = widget.NewSliderWithData(0.4, 2.0, pitchData)
	pitchSlider.Step = 0.1

	pitchLabel := canvas.NewText("Voice Pitch", theme.Color(theme.ColorNameForeground))
	leftSpacer := canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(20, 0))
	left := container.NewHBox(leftSpacer, pitchLabel)
	pitchRow := container.New(
		layout.NewBorderLayout(nil, nil, left, nil),
		left,
		pitchSlider,
	)

	LocalWindow = app.NewWindow("Read text offline")
	LocalWindow.SetOnClosed(func() {
		LocalWindowIsOpen = false
	})

	textArea = widget.NewMultiLineEntry()
	textArea.SetPlaceHolder("Enter text you want to be read")

	readButton = widget.NewButton("Read", func() {
		lines := strings.Split(textArea.Text, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			reproductor.Reproduce(line, "")

		}

	})

	controls := container.NewVBox(
		pitchRow,
		readButton,
	)

	LocalContainer = container.NewBorder(
		controls,
		nil, nil, nil,
		textArea,
	)
	LocalWindow.SetContent(LocalContainer)
	LocalWindow.Resize(fyne.NewSize(700, 500))

}

func InitLocalButton() {
	LocalButton = widget.NewButton("Read text offline", func() {
		initLocalWindow(*AppReference)
		if !LocalWindowIsOpen {
			LocalWindow.Show()
			LocalWindowIsOpen = true
		}
	})

}

package local

import (
	"image/color"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"ac-tts/internal/common"
	"ac-tts/internal/reproductor"
)

var pitchSlider *widget.Slider
var readButton *widget.Button
var textArea *widget.Entry
var saveButton *widget.Button
var filenameInput *widget.Entry
var folderButton *widget.Button
var saveDir string

var LocalContainer *fyne.Container
var LocalWindow fyne.Window
var LocalWindowIsOpen = false
var LocalButton *widget.Button
var AppReference *fyne.App
var LocalErrorWindow fyne.Window
var LocalSuccessWindow fyne.Window

const format = ".wav"

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

	filenameInput = widget.NewEntry()
	filenameInput.SetPlaceHolder("Enter file name")

	folderButton = widget.NewButton("Select Folder", func() {
		dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				initErrorWindow(app, "Error selecting folder: "+err.Error())
				LocalErrorWindow.Show()
				return
			}
			if list == nil {
				folderButton.SetText("Select Folder")
				saveDir = ""
				return
			}
			saveDir = list.Path()
			folderButton.SetText("Folder: " + saveDir)
		}, LocalWindow).Show()
	})

	saveButton = widget.NewButton("Save audio", func() {
		if filenameInput.Text == "" {
			initErrorWindow(app, "You must enter a filename")
			LocalErrorWindow.Show()
			return
		}

		if saveDir == "" {
			initErrorWindow(app, "You must select a folder")
			LocalErrorWindow.Show()
			return
		}

		lines := strings.Split(textArea.Text, "\n")
		finalLines := ""
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			finalLines = finalLines + " " + line
		}
		path := filepath.Join(saveDir, filenameInput.Text+format)
		reproductor.SaveAsWav(finalLines, path)
		initSuccessWindow(app, "Audio clip saved!")
		LocalSuccessWindow.Show()
	})

	readButton = widget.NewButton("Read", func() {
		lines := strings.Split(textArea.Text, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			reproductor.Reproduce(line, "", common.Pitch)

		}

	})

	controls := container.NewVBox(
		pitchRow,
		container.NewHBox(
			readButton,
			container.NewGridWrap(fyne.NewSize(300, filenameInput.MinSize().Height), filenameInput),
			folderButton,
			saveButton,
		),
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

func initErrorWindow(app fyne.App, msg string) {
	LocalErrorWindow = app.NewWindow("Error!")
	LocalErrorWindow.SetContent(widget.NewLabel(msg))

}

func initSuccessWindow(app fyne.App, msg string) {
	LocalSuccessWindow = app.NewWindow("Success!")
	LocalSuccessWindow.SetContent(widget.NewLabel(msg))
}

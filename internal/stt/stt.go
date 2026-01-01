package stt

import (
	"bufio"
	"context"
	"image/color"
	"os"
	"regexp"
	"strconv"
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
var initSttButton *widget.Button
var testPitchSTTButton *widget.Button
var fileButton *widget.Button

var STTContainer *fyne.Container
var STTWindow fyne.Window
var STTWindowIsOpen = false
var STTButton *widget.Button
var AppReference *fyne.App
var STTErrorWindow fyne.Window
var file *os.File
var cancel context.CancelFunc
var fileName string

func initSTTWindow(app fyne.App) {
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

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())

	STTWindow = app.NewWindow("Read from STT")
	STTWindow.SetOnClosed(func() {
		STTWindowIsOpen = false
		cancel()
	})

	testPitchSTTButton = widget.NewButton("Test Pitch", func() {
		reproductor.Reproduce("Esto es una prueba de Pitch :)", "")
	})

	fileButton = widget.NewButton("Select text file", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				initErrorWindow(app, err.Error())
				STTErrorWindow.Show()
				return
			}

			if reader == nil {
				file = nil
				fileButton.SetText("Select text file")
				return
			}

			defer reader.Close()

			file, err = os.Open(reader.URI().Path())
			if err != nil {
				STTErrorWindow.Show()
				initErrorWindow(app, "Error trying to open file: "+err.Error())
			}

			fileName = file.Name()
			fileButton.SetText(fileName)

		}, STTWindow).Show()
	})

	initSttButton = widget.NewButton("Read STT file", func() {
		if file == nil {
			initErrorWindow(app, "Select a file first")
			STTErrorWindow.Show()
			return
		}
		initSttButton.SetText("Reading...")
		initSttButton.Disable()
		go func(ctx context.Context) {
			reader := bufio.NewReader(file)
			timeRegex := regexp.MustCompile(`\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}`)
			for {

				select {
				case <-ctx.Done():
					return
				default:
					line, err := reader.ReadString('\n')
					if err != nil {
						time.Sleep(500 * time.Millisecond)
						continue
					}

					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}

					if _, err := strconv.Atoi(line); err == nil {
						continue
					}

					if timeRegex.MatchString(line) {
						continue
					}

					reproductor.Reproduce(line, "")
					time.Sleep(1500 * time.Millisecond)
				}

			}
		}(ctx)
	})

	controls := container.NewVBox(
		pitchRow,
		testPitchSTTButton,
		fileButton,
		initSttButton,
	)

	STTWindow.SetContent(controls)
	STTWindow.Resize(fyne.NewSize(700, 300))
}

func InitSTTButton() {
	STTButton = widget.NewButton("Read text from STT", func() {
		initSTTWindow(*AppReference)
		if !STTWindowIsOpen {
			STTWindow.Show()
			STTWindowIsOpen = true
		}

	})
}

func initErrorWindow(app fyne.App, msg string) {
	STTErrorWindow = app.NewWindow("Error!")
	STTErrorWindow.SetContent(widget.NewLabel(msg))

}

package main

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"github.com/faiface/beep/wav"

	"ac-tts/internal/animalese"
	"ac-tts/internal/assets"
	"ac-tts/internal/common"
	"ac-tts/internal/github"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
)

var version = "v0.0.2"
var updateTime = false

func main() {
	onlineVersion := github.GetLatestReleaseVersion()

	if onlineVersion != "" {
		updateTime = needUpdate(version, onlineVersion)
		if updateTime {
			common.InitUpdateButton()
		}
	}

	a := app.New()
	w := a.NewWindow("AC - Text to Speech for Twitch :)")

	go func() {
		web.StartWebServer()
	}()

	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})
	if err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
		panic(err)
	}

	wave := ani.AnimaleseFunc("test", true, 1.0)
	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
	}
	defer streamer.Close()
	reproductor.InitSpeaker(format)

	pitchData := binding.BindFloat(&common.Pitch)
	common.InitPitchRow(pitchData)
	common.InitTestPitchButton(func() {
		reproductor.Reproduce("Hola esto es una prueba de pitch :)", "")
	})

	common.InitConnectButton(func() { twitch.GetAuthorization() })
	w.Resize(fyne.NewSize(400, 200))

	icon := fyne.NewStaticResource("icon.png", assets.Icon)

	common.InitKofiButton()

	content := container.NewVBox(
		common.PitchRow,
		container.NewCenter(container.NewHBox(common.TestPitchButton, common.ConnectButton)),
	)
	var footer *fyne.Container
	if updateTime {
		footer = container.NewVBox(common.UpdateButton, common.KofiButton)
	} else {
		footer = container.NewVBox(common.KofiButton)
	}

	w.SetIcon(icon)
	w.SetContent(
		container.New(
			layout.NewBorderLayout(nil, footer, nil, nil),
			footer,
			container.NewStack(content),
		),
	)
	w.ShowAndRun()

}

func needUpdate(current string, online string) bool {

	currentParts := getVersionSplitted(current)
	onlineParts := getVersionSplitted(online)

	mayorCurrent, _ := strconv.Atoi(currentParts[0])
	mayorOnline, _ := strconv.Atoi(onlineParts[0])

	minorCurrent, _ := strconv.Atoi(currentParts[1])
	minorOnline, _ := strconv.Atoi(onlineParts[1])

	patchCurrent, _ := strconv.Atoi(currentParts[2])
	patchOnline, _ := strconv.Atoi(onlineParts[2])

	if mayorCurrent < mayorOnline {
		return true
	}

	if minorCurrent < minorOnline {
		return true
	}

	if patchCurrent < patchOnline {
		return true
	}

	return false

}

func getVersionSplitted(version string) []string {
	trimmed := strings.TrimPrefix(version, "v")
	return strings.Split(trimmed, ".")
}

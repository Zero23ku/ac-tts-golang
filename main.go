package main

import (
	"bytes"
	"context"
	"log"
	"os"
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
	"ac-tts/internal/tiktok"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
	"ac-tts/internal/youtube"
)

var version = "v0.3.0"
var updateTime = false

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	onlineVersion := github.GetLatestReleaseVersion()
	common.InitLeftSpacer()
	if onlineVersion != "" {
		updateTime = needUpdate(version, onlineVersion)
		if updateTime {
			common.InitUpdateButton()
		}
	}

	a := app.New()
	w := a.NewWindow("AC - Text to Speech for Twitch :) - " + version)
	w.SetOnClosed(func() {
		cancel()
		os.Exit(0)
	})
	youtube.AppReference = &a
	youtube.CTX = ctx
	youtube.InitConnectYTButton()

	tiktok.CTX = ctx
	tiktok.AppReference = &a
	tiktok.InitConnectTiktokButton()
	twitch.CTX = ctx

	go func() {
		web.StartWebServer()
	}()

	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})
	if err != nil {
		logging.CreateLog("main - Initializing animalese", err)
		log.Fatal(err)
		panic(err)
	}

	wave := ani.AnimaleseFunc("test", true, 1.0)
	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
		logging.CreateLog("main - testing Animalese", err)
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
		container.NewCenter(container.NewHBox(common.TestPitchButton, common.ConnectButton, youtube.ConnectYTButton, tiktok.ConnectTiktokButton)),
	)

	common.InitCommandCheck()
	common.InitCommandInput()

	commandContent := container.NewCenter(
		container.NewHBox(
			common.ActivateCommand,
			common.InputCommand,
		),
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
			container.NewVBox(content, commandContent),
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

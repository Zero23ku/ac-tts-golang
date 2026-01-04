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
	"ac-tts/internal/chzzk"
	"ac-tts/internal/common"
	"ac-tts/internal/github"
	"ac-tts/internal/local"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/stt"
	"ac-tts/internal/tiktok"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
	"ac-tts/internal/youtube"
)

var version = "v0.7.3"
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

	chzzk.CTX = ctx
	chzzk.AppReference = &a
	chzzk.InitConnectChzzkButton()

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
		reproductor.Reproduce("Hola esto es una prueba de pitch :)", "", common.Pitch)
	})
	common.AppReference = &a
	common.InitTwitchConnectButton(func() { twitch.GetAuthorization() })
	w.Resize(fyne.NewSize(400, 200))

	icon := fyne.NewStaticResource("icon.png", assets.Icon)

	common.InitKofiButton()

	content := container.NewVBox(
		common.PitchRow,
		container.NewCenter(container.NewHBox(common.TestPitchButton, common.TwitchConnectButton, youtube.ConnectYTButton, tiktok.ConnectTiktokButton, chzzk.ConnectChzzkButton)),
	)

	common.InitCommandCheck()
	common.InitCommandInput()
	common.InitDocsButton()
	common.InitRandomPitch()

	commandContent := container.NewCenter(
		container.NewVBox(
			container.NewHBox(
				common.ActivateCommand,
				common.InputCommand,
			),
			common.RandomPitch,
		),
	)

	local.AppReference = &a
	local.InitLocalButton()

	stt.AppReference = &a
	stt.InitSTTButton()

	var footer *fyne.Container
	if updateTime {
		footer = container.NewVBox(common.UpdateButton, common.KofiButton)
	} else {
		footer = container.NewVBox(common.DocsButton, common.KofiButton)
	}

	w.SetIcon(icon)
	w.SetContent(
		container.New(
			layout.NewBorderLayout(nil, footer, nil, nil),
			footer,
			container.NewVBox(content, commandContent, local.LocalButton, stt.STTButton),
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

package common

import (
	"fyne.io/fyne/v2/widget"
)

var ConnectButton *widget.Button

func InitConnectButton(onClick func()) {
	ConnectButton = widget.NewButton("Connect to Twitch", onClick)
}

func SetConnected() {
	ConnectButton.SetText("Connected")
	ConnectButton.Disable()
}

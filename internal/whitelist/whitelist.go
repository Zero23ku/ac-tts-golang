package whitelist

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var IsWhitelistActive = false
var whiteList *widget.Entry
var updateWhitelist *widget.Button
var whiteListCheck *widget.Check
var WhiteListContainer *fyne.Container
var UserWhitelist []string

func InitWhiteList() {
	whiteListCheck = widget.NewCheck("Active Whitelist", func(value bool) {
		IsWhitelistActive = value
	})

	whiteList = widget.NewEntry()
	whiteList.SetPlaceHolder("Enter usernames separated by ,")
	whiteList.Resize(fyne.NewSize(200, 0)) // Establecer ancho mínimo
	entryContainer := container.NewStack(whiteList)

	updateWhitelist = widget.NewButton("Update whitelist", func() {
		rawUsers := whiteList.Text
		users := strings.Split(rawUsers, ",")
		var localWhitelist []string
		for _, user := range users {
			user = strings.ToLower(strings.TrimSpace(user))
			if !IsUserInWhitelist(user) {
				localWhitelist = append(localWhitelist, user)
			}
		}
		UserWhitelist = localWhitelist
		//TODO: Solucionar
		//config.SaveConfig(!common.IsRedeemOptionActiva, common.TwitchRedeemName.Text, UserWhitelist)

	})

	WhiteListContainer = container.NewVBox(
		whiteListCheck,
		container.NewBorder(
			nil, nil, nil,
			updateWhitelist, // botón a la derecha
			entryContainer,  // entry se expande
		),
	)
}

func IsUserInWhitelist(user string) bool {
	user = strings.ToLower(strings.TrimSpace(user))

	return slices.Contains(UserWhitelist, user)
}

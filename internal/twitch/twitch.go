package twitch

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	twitch2 "github.com/joeyak/go-twitch-eventsub/v3"

	"ac-tts/internal/common"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/whitelist"
)

const TWITCH_URL = "https://id.twitch.tv/oauth2/authorize" +
	"?response_type=token" +
	"&client_id=4u4v1h8d2yfvftoqtstu0pley1pooo" +
	"&redirect_uri=http://localhost:9000" +
	"&scope=chat:read+chat:edit+channel:read:redemptions" +
	"&state=c3ab8aa609ea11e793ae92361f002671"

const TWITCH_BROADCASTER_ID = "https://api.twitch.tv/helix/users"

const CLIENT_ID = "4u4v1h8d2yfvftoqtstu0pley1pooo"

const IRC_TWITCH_SERVER = "irc.chat.twitch.tv:6667"

type RedemptionEvent struct {
	UserId               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserId    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Id                   string `json:"id"`
	UserInput            string `json:"user_input"`
	Status               string `json:"status"`
	Reward               Reward `json:"reward"`
	RedeemedAt           string `json:"redeemed_at"`
}

type Reward struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Cost   int    `json:"cost"`
	Prompt string `json:"prompt"`
}

type RedemptionCallback func(event RedemptionEvent)

var CTX context.Context

func GetAuthorization() {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", TWITCH_URL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", TWITCH_URL).Start()
	case "darwin":
		err = exec.Command("open", TWITCH_URL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		logging.CreateLog("twtich - unsupported platform", err)
		log.Fatal(err)
	}

}

var Active = false

func SubscribeToChat(token string) {
	broadcasterId, login, err := GetBroadcasterId(token)
	if err != nil {
		logging.CreateLog("twitch - couldn't get broadcaster info", err)
		log.Fatal("Error retrieving broadcaster id", err)
	}
	if common.IsRedeemOptionActiva {
		subscribeToEvent(token, broadcasterId)
		common.SetConnected()
	} else {
		conn, err := net.Dial("tcp", IRC_TWITCH_SERVER)
		if err != nil {
			logging.CreateLog("twitch - couldn't connect to twitch chat", err)
			log.Fatal("Error conectandose a IRC", err)
		}

		fmt.Fprintf(conn, "PASS %s\r\n", "oauth:"+token)
		fmt.Fprintf(conn, "NICK %s\r\n", login)
		fmt.Fprintf(conn, "JOIN #%s\r\n", login)
		common.SetConnected()
		reader := bufio.NewReader(conn)
		go func() {
			for {

				select {
				case <-CTX.Done():
					return
				default:
					line, err := reader.ReadString('\n')
					if strings.HasPrefix(line, "PING") {
						fmt.Fprintf(conn, "PONG :tmi.twitch.tv\r\n")
					}
					if err != nil {
						logging.CreateLog("twitch - couldn't get new message in chat", err)
						log.Fatal(err)
					}
					splitted := strings.Split(line, "#")
					if len(splitted) == 2 {
						message := strings.Split(splitted[1], ":")
						if len(message) == 2 {
							if Active {
								chatMsg := message[1]
								if common.IsTTSCommandActive() && strings.HasPrefix(chatMsg, common.GetTTSCommand()) {
									if common.IsPitchRandom {
										if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(message[0]) {
											reproductor.Reproduce(chatMsg, message[0], common.GetRandomPitch(), common.IsPitchRandom)
										} else if !whitelist.IsWhitelistActive {
											reproductor.Reproduce(chatMsg, message[0], common.GetRandomPitch(), common.IsPitchRandom)
										}

									} else {

										if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(message[0]) {
											reproductor.Reproduce(chatMsg, message[0], common.Pitch, common.IsPitchRandom)
										} else if !whitelist.IsWhitelistActive {
											reproductor.Reproduce(chatMsg, message[0], common.Pitch, common.IsPitchRandom)
										}
									}
								} else if !common.IsTTSCommandActive() {
									if common.IsPitchRandom {
										if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(message[0]) {
											reproductor.Reproduce(chatMsg, message[0], common.GetRandomPitch(), common.IsPitchRandom)
										} else if !whitelist.IsWhitelistActive {
											reproductor.Reproduce(chatMsg, message[0], common.GetRandomPitch(), common.IsPitchRandom)
										}
									} else {
										if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(message[0]) {
											reproductor.Reproduce(chatMsg, message[0], common.Pitch, common.IsPitchRandom)
										} else if !whitelist.IsWhitelistActive {
											reproductor.Reproduce(chatMsg, message[0], common.Pitch, common.IsPitchRandom)
										}
									}
								}

							} else if strings.Compare(strings.TrimSpace(message[1]), "End of /NAMES list") == 0 && !Active {
								Active = true
							}

						}
					}
				}

			}
		}()

	}

}

func GetBroadcasterId(token string) (string, string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", TWITCH_BROADCASTER_ID, nil)
	if err != nil {
		logging.CreateLog("twitch - couldn't create HTTP request", err)
		log.Fatal("Error creating request", err)
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", CLIENT_ID)

	resp, err := client.Do(req)
	if err != nil {
		logging.CreateLog("twitch - couldn't make HTTP request", err)
		log.Fatal("Error sending request", err)
		return "", "", err
	}
	defer resp.Body.Close()

	var result common.Response

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logging.CreateLog("twitch - couldn't deserealize response", err)
		log.Fatal("Error decoding JSON response", err)
		return "", "", err
	}

	if len(result.Data) > 0 {
		b := result.Data[0]
		broadcasterid := b.Id
		login := b.Login
		return broadcasterid, login, nil
	}

	return "", "", errors.New("No data in")
}

// Source - https://stackoverflow.com/a
// Posted by user5728991, modified by community. See post 'Timeline' for change history
// Retrieved 2025-12-06, License - CC BY-SA 4.0

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func subscribeToEvent(accessToken string, userID string) {
	client := twitch2.NewClient()
	client.OnError(func(err error) {
		logging.CreateLog("twitch - On Error Subscribe to event - ", err)
		log.Fatal("twitch - On Error Subscribe to even", err)
	})
	client.OnWelcome(func(message twitch2.WelcomeMessage) {

		events := []twitch2.EventSubscription{
			"channel.channel_points_custom_reward_redemption.add",
		}

		for _, event := range events {
			_, err := twitch2.SubscribeEvent(twitch2.SubscribeRequest{
				SessionID:   message.Payload.Session.ID,
				ClientID:    CLIENT_ID,
				AccessToken: accessToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": userID,
				},
			})
			if err != nil {
				logging.CreateLog("twitch - Event error - ", err)
				return
			}
		}
	})

	client.OnRawEvent(func(event string, metadata twitch2.MessageMetadata, subscription twitch2.PayloadSubscription) {
		//fmt.Printf("EVENT[%s]: %s: %s\n", subscription.Type, metadata, event)
		var r RedemptionEvent
		if err := json.Unmarshal([]byte(event), &r); err != nil {
			panic(err)
		}
		if common.TwitchRedeemName.Text == r.Reward.Title {
			if common.IsPitchRandom {
				reproductor.Reproduce(r.UserInput, r.UserName, common.GetRandomPitch(), common.IsPitchRandom)
			} else {
				reproductor.Reproduce(r.UserInput, r.UserName, common.Pitch, common.IsPitchRandom)
			}
		}
	})

	go func() {
		if err := client.Connect(); err != nil {
			logging.CreateLog("twitch - Could not connect client - ", err)
			log.Fatal("twitch - Could not connect client - ", err)
		}
	}()

	go func() {
		<-CTX.Done()
		client.Close()
	}()
}

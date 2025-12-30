# Animal Crossing TTS for Streamers

Here's gonna explain how to use this app.

This is the main app's window. Here you can connect to different stream services and also test the pitch of your TTS.

To change the pitch, you just need to move the slider, then click "Test Voice" to hear how it will sound.

![AC TTS](./img/app.png)

## Disclaimer

This application is totally FREE and Open Source, there is no monetization in any shape or form and never will be, you are free to use, transform, or do anything with this code as long as you follow the GPLv3 license

If you want to support me, you can do it by sending me a Ko-fi

[![Ko-fi](./img/kofi4.png)](https://ko-fi.com/harpuia)

## Connect to Twitch

When you click "Connecto to Twitch", the following window will appears in your browser, asking you to give permission to read and write messages in your chat (At the moment, this app don't send any message to chat, it may change in the future if there is a interested feature to implement, or it will be removed in the future)

![Twitch's window](./img/twitch-access.png)

After you authorize the app, a new tab pointing to http://localhost:9000 should open and close by itself, if it doesn't close by itself, you can close it without issue. This is done because the way Twitch send it's authorization token can't be catch in any other way other than having a little frontend webapp to read the url and send it from there to the TTS backend, this is all being done locally.

If everything went right, the "Connect to Twitch" button should appears like this

![You are connected to Twitch](./img/twitch-connected.png)

## Connect to Tiktok

When you click the "Connect to Tiktok" button, the following window will appears

![Tiktok's window](./img/tiktok-window.png)

To connect to your Tiktok live chat, first you need to be already livestreaming and then, insert your Tiktok username in the field shown there. It may launch and error and ask you to try again, so you may have to try few times, but if Tiktok's window closes without showing any error, then it is connected.

## Connect to Youtube

Due to how Youtube's chat works, this implementation is way more complicated than the other platforms, I'm still looking for options, so this may or may not change in the future.

To connect to your Youtube's livestream chat, you need to enter your livestream url (Ex: https://youtube.com/watch?v=rnXIjl_Rzy4) and your Youtube's API key.

Due to Youtube's Api key limits and quota, for now, the best way to ensure this works for everyone is that **every user need to bring they own API Key**

To leanr how to obtain your own Youtube's API key you can [follow official guide](https://developers.google.com/youtube/registering_an_application) or [check this video](https://www.youtube.com/watch?v=EPeDTRNKAVo)

**Please ve careful, don't share your API Key, treat it as any other of your passwords, this app doesn't save your API Key so you will need to use it every time you want to connect it to your chat**

![Youtube's window](./img/youtube-window.png)
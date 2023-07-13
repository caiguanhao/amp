# amp

Backup and restore Apple Music Playlist.

## Usage

Put token in config file ~/.amp.json:

```
{
  "developer_token": "eyJhb...",
  "user_token": "AiGUK..."
}
```

Backup playlists to current directory:

```
amp backup

2023/07/13 14:17:16 found 2 playlists
2023/07/13 14:17:17 writing playlist FAVORITES to playlist.p.EYWrg35c4DOAgv.json
2023/07/13 14:17:18 writing playlist Replay 2022 to playlist.p.B0A8vz3h0A58RK.json
```

Restore specific playlist:

```
amp restore playlist.p.EYWrg35c4DOAgv.json TEST

2023/07/13 14:17:43 created playlist [ name: TEST , id: p.dl0k4dDhWZkrEN ] with 17 tracks
```

## Token

1. Copy the Team ID in [Apple Developer Account](https://developer.apple.com/account) page.
2. [Create](https://developer.apple.com/account/resources/identifiers/list/musicId) MusicKit identifier.
3. [Create](https://developer.apple.com/account/resources/authkeys/list) and download MusicKit key file, copy the key ID.
4. [Generate](https://github.com/minchao/go-apple-music/blob/master/examples/token-generator/main.go) developer token with team ID, key file and key ID.
5. In any Chrome page, use DevTools console. Enter following script, log in and copy the user token:
```js
var script = document.createElement('script')
script.src = 'https://js-cdn.music.apple.com/musickit/v1/musickit.js'
script.onload = function () {
  const button = document.createElement('button')
  button.addEventListener('click', function () {
    MusicKit.configure({
      developerToken: 'YOUR_DEVELOPER_TOKEN',
    }).authorize().then(function (token) {
      console.log(token)
    })
  })
  button.click()
}
document.head.appendChild(script)
```

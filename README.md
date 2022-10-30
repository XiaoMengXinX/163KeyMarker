# 163 Key Marker

[![Go Report Card](https://goreportcard.com/badge/github.com/XiaoMengXinX/163KeyMarker)](https://goreportcard.com/report/github.com/XiaoMengXinX/163KeyMarker)
[![](https://pkg.go.dev/badge/github.com/XiaoMengXinX/163KeyMarker)](https://pkg.go.dev/github.com/XiaoMengXinX/163KeyMarker)

A simple tool for adding id3v2 tag and [163 key](https://stageguard.top/2019/10/27/analyze-163-music-key) tag to a music
file.

## Some examples:

- Add 163 key, id3v2 tag and cover to a music file:

```go
package main

import (
	"net/http"
	"os"

	"github.com/XiaoMengXinX/163KeyMarker"
	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
)

func main() {
	data := utils.RequestData{}

	songDetail, _ := api.GetSongDetail(data, []int{1965687934})
	songUrlData, _ := api.GetSongURL(data, api.SongURLConfig{Ids: []int{1965687934}})

	pic, _ := http.Get(songDetail.Songs[0].Al.PicUrl)

	mark := marker.CreateMarker(songDetail.Songs[0], songUrlData.Data[0])

	file, _ := os.Open("1965687934.flac")
	defer file.Close()

	err := marker.AddMusicID3V2(file, pic.Body, mark)
	if err != nil {
		panic(err)
	}
}
```

The output file will look like:

![Snipaste_2022-10-30_21-30-41](https://user-images.githubusercontent.com/19994286/198881269-8e4f1a41-b277-4c58-82bf-95d7ff16f888.png)

- Parse the embedded 163 key tag from a music file:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/XiaoMengXinX/163KeyMarker"
)

func main() {
	file, _ := os.Open("1965687934.flac")
	defer file.Close()

	mark, err := marker.Parse163Key(file)
	if err != nil {
		panic(err)
	}

	markerJson, _ := json.MarshalIndent(mark, "", "  ")

	fmt.Println(string(markerJson))
}
```

Output:

```json
{
  "musicId": 1965687934,
  "musicName": "洛天依印象曲·华风夏韵",
  "artist": [
    [
      "Tino.S3",
      12002114
    ]
  ],
  "albumId": 148413696,
  "album": "洛天依印象曲·华风夏韵",
  "albumPicDocId": "109951167697837238",
  "albumPic": "https://p1.music.126.net/1oZ0ftxNE5DGzwvO97bfLA==/109951167697837238.jpg",
  "bitrate": 320000,
  "mp3DocId": "7ee41abd08e7e829b41dfd7b0689a256",
  "duration": 207158,
  "mvId": 0,
  "alias": [],
  "format": "mp3"
}
```

- Create a 163 key string from marker data:

```go
package main

import (
	"fmt"

	marker "github.com/XiaoMengXinX/163KeyMarker"
)

func main() {
	data := marker.MarkerData{
		MusicId:   1965687934,
		MusicName: "洛天依印象曲·华风夏韵",
	}
	mark := marker.Create163KeyStr(data)
	fmt.Println(mark)
}
```

Output:

```
163 key(Don't modify):L64FU3W4YxX3ZFTmbZ+8/ZedOKM0To3lAHr1q3yMEARRJ+Dh+02XUDlnBUo9rPC7jWKxG/Gy7ZpnrH1ckjYKYK+JIykH3KOIxteIZFAWNkWJfR2PtFQs5GIQ921WIUnZai33f5lhpDb3hPlVWOzxPRg136s014agaLb9aILz/o7nRBknv1hFWdGpvpMn3vk2vkQy6ExluHEbeU7tdPZ/ENE44YGEWHTu/DkkFJkQ1OrPmcTkHsEAQqBMQ0pBox7lh3Nh3Iib5LNuGr3vENJ21SMwCerf0QUQzSobspcWLzc=
```

- Decrypt the 163 key string:

```go
package main

import (
	"fmt"

	"github.com/XiaoMengXinX/163KeyMarker"
)

func main() {
	encrypted := "163 key(Don't modify):L64FU3W4YxX3ZFTmbZ+8/ZedOKM0To3lAHr1q3yMEARRJ+Dh+02XUDlnBUo9rPC7jWKxG/Gy7ZpnrH1ckjYKYK+JIykH3KOIxteIZFAWNkXQFFl4yALfP5TWmM7sYeBMRLPdbSFzK6VaPrsN9FTCJ8IAo1F/rUBzF2OnlxZrQuwSQXmcPESzQbbgOeDv9V58u+0v3r/fpEQkFcFeIjsR1kK6EGutL/WX6tWp2t6Zfx3aIJLt/ZOLKeUhlvYGxTtNDWqaqg19i0CKZ4Lpay5Ha1GDqBu9VZ7ZPKo2ofJ0cLaPuQqEVjESFiQQPdhrj6hMSFo+rZpI4WKkZ6xwR6vQHmXnlDXh3NzT/vifHNoTBCiZM3PGicKQQKq4KsRwU1Gd2L12na0AeVuxkhHpKAIYo0/eOOrDfcOqh+d2xkdrWEOCGx06KlZTOOJk1X6lQjNUbp5G4VmKySCPxkPdCf/mUbAns8+yfb2WJgxUeJ40glqR8v0JixTjqN7ssZ1yegpkn1I8rS+tkrpbegWZeTOf7Kywi4SypxwQOc9S/+QUi9s="
	decrypted := marker.Decrypt163key(encrypted)
	fmt.Println(decrypted)
}

```

Output:

```
music:{"musicId": 1965687934, "musicName": "洛天依印象曲·华风夏韵", "artist":[["Tino.S3", 12002114]], "albumId": 148413696, "album": "洛天依印象曲·华风夏韵", "albumPicDocId": "109951167697837238", "albumPic": "https://p1.music.126.net/1oZ0ftxNE5DGzwvO97bfLA==/109951167697837238.jpg", "bitrate": 320000, "mp3DocId": "7ee41abd08e7e829b41dfd7b0689a256", "duration": 207158, "mvId": 0, "alias":[], "format": "mp3"}
```
  

package marker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
)

// CreateMarker creates marker data from SongDetailData and SongURLData
func CreateMarker(songDetail types.SongDetailData, songUrl types.SongURLData) (marker MarkerData) {
	var artists [][]interface{}
	for _, j := range songDetail.Ar {
		var artist []interface{}
		artist = make([]interface{}, 2)
		artist[0] = j.Name
		artist[1] = j.Id
		artists = append(artists, artist)
	}
	return MarkerData{
		MusicId:       songDetail.Id,
		MusicName:     songDetail.Name,
		Artist:        artists,
		AlbumId:       songDetail.Al.Id,
		Album:         songDetail.Al.Name,
		AlbumPicDocId: songDetail.Al.PicStr,
		AlbumPic:      songDetail.Al.PicUrl,
		Bitrate:       songUrl.Br,
		Mp3DocId:      songUrl.Md5,
		Duration:      songDetail.Dt,
		MvId:          songDetail.Mv,
		Alias:         songDetail.Alia,
		Format:        songUrl.Type,
	}
}

// Create163KeyStr formats the marker data to 163 key string
func Create163KeyStr(marker MarkerData) (markerText string) {
	markerJson, err := json.Marshal(marker)
	if err != nil {
		return markerText
	}
	decryptedMarker := base64.StdEncoding.EncodeToString(utils.MarkerEncrypt(fmt.Sprintf("music:%s", string(markerJson))))
	markerText = fmt.Sprintf("163 key(Don't modify):%s", decryptedMarker)
	return markerText
}

// Encrypt163Key encrypts the 163 key string
func Encrypt163Key(decrypted string) (encrypted string) {
	return base64.StdEncoding.EncodeToString(utils.MarkerEncrypt(decrypted))
}

// Decrypt163key decrypts the 163 key string
func Decrypt163key(encrypted string) (decrypted string) {
	strings.TrimPrefix(encrypted, "163 key(Don't modify):")
	data, _ := base64.StdEncoding.DecodeString(encrypted)
	return string(utils.MarkerDecrypt(data))
}

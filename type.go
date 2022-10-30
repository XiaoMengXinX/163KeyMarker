package marker

// MarkerData is the format data of 163 key
type MarkerData struct {
	MusicId       int             `json:"musicId"`
	MusicName     string          `json:"musicName"`
	Artist        [][]interface{} `json:"artist"`
	AlbumId       int             `json:"albumId"`
	Album         string          `json:"album"`
	AlbumPicDocId string          `json:"albumPicDocId"`
	AlbumPic      string          `json:"albumPic"`
	Bitrate       int             `json:"bitrate"`
	Mp3DocId      string          `json:"mp3DocId"`
	Duration      int             `json:"duration"`
	MvId          int             `json:"mvId"`
	Alias         []interface{}   `json:"alias"`
	Format        string          `json:"format"`
}

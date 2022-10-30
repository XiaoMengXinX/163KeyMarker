package marker

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bogem/id3v2"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

// AddMusicID3V2 detects the music file type and adds ID3V2 tags to it.
func AddMusicID3V2(musicFile, picFile io.Reader, marker MarkerData) (err error) {
	file, ok := musicFile.(*os.File)
	if !ok {
		return fmt.Errorf("musicFile was not initialized with file")
	}
	fileType := detectFileType(file)
	switch fileType {
	case FileTypeMp3:
		return AddMp3Id3v2(file, picFile, marker)
	case FileTypeFlac:
		return AddFlacId3v2(file, picFile, marker)
	}
	return fmt.Errorf("invaid file type")
}

// AddMp3Id3v2 adds ID3V2 tags to mp3 file.
func AddMp3Id3v2(musicFile *os.File, picFile io.Reader, marker MarkerData) (err error) {
	musicTag, _ := id3v2.ParseReader(musicFile, id3v2.Options{Parse: false})
	defer musicTag.Close()
	musicTag.SetDefaultEncoding(id3v2.EncodingUTF8)
	musicTag.SetTitle(marker.MusicName)
	musicTag.SetArtist(FormatArtistsStr(marker))
	if marker.Album != "" {
		musicTag.SetAlbum(marker.Album)
	}
	comment := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingISO,
		Language:    "chs",
		Description: "",
		Text:        Create163KeyStr(marker),
	}
	musicTag.AddCommentFrame(comment)
	if picFile != nil {
		artwork, err := io.ReadAll(picFile)
		if err != nil {
			return fmt.Errorf("error while reading album pic: %v ", err)
		}
		mime := http.DetectContentType(artwork[:32])
		pic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingISO,
			MimeType:    mime,
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     artwork,
		}
		musicTag.AddAttachedPicture(pic)
	}
	if err := musicTag.Save(); err != nil {
		return fmt.Errorf("error while saving: %v ", err)
	}
	return nil
}

// AddFlacId3v2 adds ID3V2 tags to flac file.
func AddFlacId3v2(musicFile *os.File, picFile io.Reader, marker MarkerData) (err error) {
	var file flacFile
	err = file.Parse(musicFile)
	if err != nil {
		return err
	}
	tag := flacvorbis.New()

	if picFile != nil {
		artwork, err := io.ReadAll(picFile)
		if err != nil {
			return fmt.Errorf("error while reading album pic: %v ", err)
		}
		mime := http.DetectContentType(artwork[:32])
		picture, err := flacpicture.NewFromImageData(flacpicture.PictureTypeFrontCover, "Front cover", artwork, mime)
		if err == nil {
			pictureMeta := picture.Marshal()
			file.file.Meta = append(file.file.Meta, &pictureMeta)
		}

	}

	_ = tag.Add(flacvorbis.FIELD_TITLE, marker.MusicName)
	_ = tag.Add(flacvorbis.FIELD_ARTIST, FormatArtistsStr(marker))
	if marker.Album != "" {
		_ = tag.Add(flacvorbis.FIELD_ALBUM, marker.Album)
	}
	_ = tag.Add(flacvorbis.FIELD_DESCRIPTION, Create163KeyStr(marker))

	tagMeta := tag.Marshal()

	var idx int
	for i, m := range file.file.Meta {
		if m.Type == flac.VorbisComment {
			idx = i
			break
		}
	}
	if idx > 0 {
		file.file.Meta[idx] = &tagMeta
	} else {
		file.file.Meta = append(file.file.Meta, &tagMeta)
	}

	if err = file.Save(); err != nil {
		return fmt.Errorf("error while saving: %v ", err)
	}

	return err
}

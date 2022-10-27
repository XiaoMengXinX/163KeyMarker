package marker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

// Parse163Key parse 163 key from reader
func Parse163Key(file *os.File) (marker MarkerData, err error) {
	fileType := detectFileType(file)
	switch fileType {
	case FileTypeMp3:
		return ReadMp3Key(file)
	case FileTypeFlac:
		return ReadFlacKey(file)
	}
	return marker, fmt.Errorf("invaid file type ")
}

func ReadMp3Key(file io.Reader) (marker MarkerData, err error) {
	tag, err := id3v2.ParseReader(file, id3v2.Options{Parse: true})
	if err != nil {
		return marker, err
	}
	defer tag.Close()

	var comment string
	frames := tag.GetFrames(tag.CommonID("Comments"))
	if len(frames) != 0 {
		val, ok := frames[0].(id3v2.CommentFrame)
		if !ok {
			return marker, fmt.Errorf("couldn't assert comment frame ")
		}
		comment = val.Text
	}

	if strings.Contains(comment, "163 key(Don't modify):") {
		markerText := strings.TrimPrefix(comment, "163 key(Don't modify):")
		markerJson := strings.Replace(Decrypt163key(markerText), "music:", "", 1)
		var marker MarkerData
		_ = json.Unmarshal([]byte(markerJson), &marker)
		return marker, err
	}

	return marker, fmt.Errorf("invaid comment frame ")
}

func ReadFlacKey(file io.Reader) (marker MarkerData, err error) {
	flacFile, err := flac.ParseMetadata(file)
	if err != nil {
		return marker, err
	}

	var tag *flacvorbis.MetaDataBlockVorbisComment
	for _, meta := range flacFile.Meta {
		if meta.Type == flac.VorbisComment {
			tag, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			if err != nil {
				panic(err)
			}
		}
	}

	comment, err := tag.Get("DESCRIPTION")
	if err != nil {
		return marker, err
	}
	if strings.Contains(comment[0], "163 key(Don't modify):") && len(comment) != 0 {
		markerText := strings.TrimPrefix(comment[0], "163 key(Don't modify):")
		markerJson := strings.Replace(Decrypt163key(markerText), "music:", "", 1)
		var marker MarkerData
		_ = json.Unmarshal([]byte(markerJson), &marker)
		return marker, err
	}

	return marker, fmt.Errorf("invaid Comment Frame ")
}

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

// Parse163Key detects file type and parses the embed 163 key from file
func Parse163Key(musicFile io.Reader) (marker MarkerData, err error) {
	file, ok := musicFile.(*os.File)
	if !ok {
		return marker, fmt.Errorf("musicFile was not initialized with file")
	}
	fileType := detectFileType(file)
	switch fileType {
	case FileTypeMp3:
		return Parse163KeyFromMp3File(file)
	case FileTypeFlac:
		return Parse163KeyFromFlacFile(file)
	}
	return marker, fmt.Errorf("invaid file type")
}

// Parse163KeyFromMp3File parses the embed 163 key from mp3 file
func Parse163KeyFromMp3File(file *os.File) (marker MarkerData, err error) {
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

// Parse163KeyFromFlacFile parses the embed 163 key from flac file
func Parse163KeyFromFlacFile(file *os.File) (marker MarkerData, err error) {
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
		markerJson := strings.TrimPrefix(Decrypt163key(markerText), "music:")
		var marker MarkerData
		_ = json.Unmarshal([]byte(markerJson), &marker)
		return marker, err
	}

	return marker, fmt.Errorf("invaid Comment Frame ")
}

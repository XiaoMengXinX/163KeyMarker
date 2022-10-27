package marker

import (
	"fmt"
	"os"
	"strings"
)

const (
	FileTypeMp3 = iota
	FileTypeFlac
)

func detectFileType(file *os.File) int {
	buf := make([]byte, 32)
	n, _ := file.Read(buf)

	_, _ = file.Seek(0, 0)

	fileCode := fmt.Sprintf("%X", buf[:n])

	if strings.HasPrefix(fileCode, "494433") || strings.HasPrefix(fileCode, "FFFB") {
		return FileTypeMp3
	} else if strings.HasPrefix(fileCode, "664C6143") {
		return FileTypeFlac
	}
	return -1
}

func formatArtistsStr(marker MarkerData) string {
	var artists string
	for i, ar := range marker.Artist {
		if i == 0 {
			artists = ar[0].(string)
		} else {
			artists = fmt.Sprintf("%s, %s", artists, ar[0].(string))
		}
	}
	return artists
}

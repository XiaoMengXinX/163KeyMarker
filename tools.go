package marker

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	FileTypeMp3 = iota
	FileTypeFlac
)

// FormatArtistsStr formats the artists slice into a string.
// For example, if the artists slice is ["A", "B", "C"], the result will be "A, B, C".
func FormatArtistsStr(marker MarkerData) string {
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

var bsPool = sync.Pool{
	New: func() interface{} { return nil },
}

func getByteSlice(size int) []byte {
	fromPool := bsPool.Get()
	if fromPool == nil {
		return make([]byte, size)
	}
	bs := fromPool.([]byte)
	if cap(bs) < size {
		bs = make([]byte, size)
	}
	return bs[0:size]
}

func putByteSlice(b []byte) {
	bsPool.Put(b)
}

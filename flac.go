package marker

import (
	"fmt"
	"io"
	"os"

	"github.com/go-flac/go-flac"
)

type flacFile struct {
	file   *flac.File
	reader io.Reader
}

// Parse parses the flac file and stores the metadata in the flacFile struct.
func (f *flacFile) Parse(r io.Reader) (err error) {
	f.file, err = flac.ParseMetadata(r)
	f.reader = r
	return err
}

// Save saves the flac file with metadata.
func (f *flacFile) Save() (err error) {
	file, ok := f.reader.(*os.File)
	if !ok {
		return fmt.Errorf("flacFile was not initialized with file")
	}

	originalStat, err := file.Stat()
	if err != nil {
		return err
	}

	name := file.Name() + "-id3v2"
	newFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, originalStat.Mode())
	if err != nil {
		return err
	}

	tempfileShouldBeRemoved := true
	defer func() {
		if tempfileShouldBeRemoved {
			_ = os.Remove(newFile.Name())
		}
	}()

	_, _ = newFile.Write([]byte("fLaC"))
	for i, meta := range f.file.Meta {
		last := i == len(f.file.Meta)-1
		_, _ = newFile.Write(meta.Marshal(last))
	}

	buf := getByteSlice(128 * 1024)
	defer putByteSlice(buf)
	if _, err = io.CopyBuffer(newFile, file, buf); err != nil {
		return err
	}

	_ = newFile.Close()
	_ = file.Close()

	if err = os.Rename(newFile.Name(), file.Name()); err != nil {
		return err
	}
	tempfileShouldBeRemoved = false

	f.reader, err = os.Open(file.Name())
	if err != nil {
		return err
	}

	return nil
}

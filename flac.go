package marker

import (
	"io"
	"os"

	"github.com/go-flac/go-flac"
)

type flacFile struct {
	file   *flac.File
	reader *os.File
}

func (f *flacFile) Parse(r *os.File) (err error) {
	f.file, err = flac.ParseMetadata(r)
	f.reader = r
	return err
}

func (f *flacFile) Save() (err error) {
	originalStat, err := f.reader.Stat()
	if err != nil {
		return err
	}

	name := f.reader.Name() + "-id3v2"
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

	buf := make([]byte, 5000)
	for {
		n, err := io.ReadFull(f.reader, buf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := newFile.Write(buf[:n]); err != nil {
			return err
		}
	}

	_ = newFile.Close()
	_ = f.reader.Close()

	if err = os.Rename(newFile.Name(), f.reader.Name()); err != nil {
		return err
	}
	tempfileShouldBeRemoved = false

	f.reader, err = os.Open(f.reader.Name())
	if err != nil {
		return err
	}

	return nil
}

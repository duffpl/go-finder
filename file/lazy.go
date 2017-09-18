package file

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type ChecksumCallback func(path string) ([]byte, error)
type MimeCallback func(path string) (string, error)

type lazyFileInfo struct {
	os.FileInfo

	abs      string
	mime     string
	checksum []byte

	mimeCallback     MimeCallback
	checksumCallback ChecksumCallback
}

func (f *lazyFileInfo) Mime() (result string, err error) {
	if f.mime == "" {
		f.mime, err = f.mimeCallback(f.abs)
		if err != nil {
			err = errors.Wrap(err, "mime")
			return
		}
	}
	result = f.mime
	return
}

func (f *lazyFileInfo) Abs() (result string, err error) {
	return f.abs, nil
}

func (f *lazyFileInfo) Checksum() (result []byte, err error) {
	if f.checksum == nil {
		f.checksum, err = f.checksumCallback(f.abs)
		if err != nil {
			err = errors.Wrap(err, "checksum")
			return
		}
	}
	result = f.checksum
	return
}

// NewLazyFileInfoExByPath creates new lazyFileInfo instance
func NewLazyFileInfoExByPath(path string, csCb ChecksumCallback, mCb MimeCallback) (result FileInfoEx, err error) {
	var (
		stat os.FileInfo
		abs  string
	)
	if stat, err = os.Stat(path); err != nil {
		err = errors.Wrap(err, "os.Stat")
		return
	}
	if abs, err = filepath.Abs(path); err != nil {
		err = errors.Wrap(err, "filepath.Abs")
		return
	}
	result = &lazyFileInfo{
		FileInfo:         stat,
		abs:              abs,
		mimeCallback:     mCb,
		checksumCallback: csCb,
	}
	return
}

package os

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type checksumCallback func(path string) ([]byte, error)
type mimeCallback func(path string) (string, error)

type lazyFileInfo struct {
	os.FileInfo

	abs      string
	mime     string
	checksum []byte

	mimeCallback     mimeCallback
	checksumCallback checksumCallback
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

func NewCollectionFromPaths(paths []string, csCb checksumCallback, mCb mimeCallback) (result []FileInfoEx, err error) {
	for _, path := range paths {
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
		result = append(result, &lazyFileInfo{
			FileInfo:         stat,
			abs:              abs,
			mimeCallback:     mCb,
			checksumCallback: csCb,
		})
	}
	return
}

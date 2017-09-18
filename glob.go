package finder

import (
	"github.com/pkg/errors"
	"github.com/duffpl/go-finder/file"
)

type GlobFunc func(pattern string) ([]string, error)
type FileInfoExGlobFunc func(pattern string) ([]file.FileInfoEx, error)

// NewLazyGlobber creates function that uses result of fileinfo.Glob to create slice of file.FileInfoEx items with
// injected checksum and mime callbacks
func NewLazyGlobber(gf GlobFunc, csCb file.ChecksumCallback, mCb file.MimeCallback) FileInfoExGlobFunc {
	return func(pattern string) (result []file.FileInfoEx, err error) {
		var matches []string
		if matches, err = gf(pattern); err != nil {
			err = errors.Wrap(err, "glob")
			return
		}
		var info file.FileInfoEx
		for _, match := range matches {
			if info, err = file.NewLazyFileInfoExByPath(match, csCb, mCb); err != nil {
				err = errors.Wrap(err, "new fileinfoex")
				return
			}
			result = append(result, info)
		}
		return
	}
}
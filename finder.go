package finder

import (
	"os"
	"path/filepath"
	"github.com/pkg/errors"
	"fmt"
	"github.com/duffpl/go-finder/mimechecker"
)

type Finder struct {
	filters []func(*FileInfoEx) (bool, error)
	err     error
}

var (
	stdGlobber globber
	stdMimeChecker mimechecker.Checker
)

func init() {
	stdGlobber = &doubleStarGlobber{}
	stdMimeChecker = mimechecker.NewMulti(mimechecker.NewGoHttp(), mimechecker.NewGoMime())
}

func New() *Finder {
	return &Finder{}
}

func (f *Finder) Glob(glob string, relativeTo string) (result []*FileInfoEx, err error) {
	globResult, err := stdGlobber.Glob(glob)
	fmt.Println(err)
	if err != nil {
		err = errors.Wrap(err, "Glob:stdGlobber")
		return
	}
	for _, globItem := range globResult {
		fileInfo, infoErr := newFileInfoEx(globItem, relativeTo)
		if infoErr != nil {
			err = errors.Wrap(infoErr, "Glob:process")
			return
		}
		if f.checkFilters(fileInfo) {
			result = append(result, fileInfo)
		}
	}
	return
}

func (f *Finder) checkFilters(input *FileInfoEx) bool {
	for _, filterFunction := range f.filters {
		matched, _ := filterFunction(input)
		if !matched {
			return false
		}
	}
	return true
}

func (f *Finder) addFilter(filter func(fiex *FileInfoEx) (bool, error)) {
	f.filters = append(f.filters, filter)
}

func newFileInfoEx(path string, relTo string) (result *FileInfoEx, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "cannot create FileInfoEx")
		}
	}()
	stat, err := os.Stat(path)
	if err != nil {
		err = errors.Wrap(err, "os.Stat")
		return
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		err = errors.Wrap(err, "filepath.Abs")
		return
	}
	rel, err := filepath.Rel(relTo, abs)
	if err != nil {
		err = errors.Wrap(err, "filepath.Rel")
		return
	}
	result = &FileInfoEx{
		FileInfo:     stat,
		AbsolutePath: path,
		RelativePath: rel,
	}
	return
}

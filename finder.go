package finder

import (
	"github.com/pkg/errors"
	"github.com/duffpl/go-finder/mimechecker"
	"github.com/duffpl/go-finder/os"
	"sort"
	"github.com/bmatcuk/doublestar"
	"crypto/md5"
	os2 "os"
	"io"
	"sync"
)


type GlobFunc func(pattern string) (result []os.FileInfoEx, err error)

type Finder struct {
	numCheckers int
	globFunc GlobFunc
	filters []struct {
		callback func(ex os.FileInfoEx) (bool, error)
		order  int
	}
	err error
}

var(
	defaultGlobFunc func(pattern string) (result []os.FileInfoEx, err error)
	defaultChecksumCallback func(path string) (result []byte, err error)
	defaultMimeCallback func(path string) (result string, err error)
	defaultMimeChecker mimechecker.Checker
	defaultCheckersConcurrency int = 8
)

func MD5ByPath(path string) (result []byte, err error) {
	checksumWriter := md5.New()
	handle, err := os2.Open(path)
	if err != nil {
		err = errors.Wrap(err, "os.Open")
		return
	}
	_, err = io.Copy(checksumWriter, handle)
	if err != nil {
		err = errors.Wrap(err, "io.Copy")
		return
	}
	defer handle.Close()
	result = checksumWriter.Sum(nil)
	return
}

func init() {
	defaultMimeChecker = mimechecker.NewMulti(mimechecker.NewGoHttp(), mimechecker.NewGoMime())
	defaultChecksumCallback = MD5ByPath
	defaultMimeCallback = func(path string) (result string, err error) {
		return defaultMimeChecker.TypeByFile(path)
	}
	defaultGlobFunc = func(pattern string) (result []os.FileInfoEx, err error) {
		paths, err := doublestar.Glob(pattern)
		if err != nil {
			err = errors.Wrap(err, "glob")
			return
		}
		return os.NewCollectionFromPaths(paths, defaultChecksumCallback, defaultMimeCallback)
	}
}

func New() *Finder {
	return new(Finder).
	SetGlobFunc(defaultGlobFunc).
	SetCheckerConcurrency(defaultCheckersConcurrency)
}

func (f *Finder) SetGlobFunc(gf GlobFunc) *Finder {
	f.globFunc = gf
	return f
}

func (f *Finder) SetCheckerConcurrency(num int) *Finder {
	if num <= 0 {
		f.err = errors.New("checkers count must be larger than 0")
	} else {
		f.numCheckers = num
	}
	return f
}

func (f *Finder) Glob(pattern string) (result []os.FileInfoEx, err error) {
	if f.err != nil {
		err = f.err
		return
	}
	globResult, err := f.globFunc(pattern)
	if err != nil {
		err = errors.Wrap(err, "glob")
		return
	}
	checkersInput := make(chan os.FileInfoEx)
	checkersOutput := make(chan os.FileInfoEx)
	wg := &sync.WaitGroup{}

	wg.Add(len(globResult))
	for i:=0; i<f.numCheckers; i++ {
		go func() {
			for fex := range checkersInput {
				if f.checkFilters(fex) {
					checkersOutput<-fex
				} else {
					wg.Done()
				}
			}
		}()
	}

	go func() {
		for fex := range checkersOutput {
			result = append(result, fex)
			wg.Done()
		}
	}()
	for _, globItem := range globResult {
		checkersInput<-globItem
	}
	close(checkersInput)
	wg.Wait()
	close(checkersOutput)
	return
}

func (f *Finder) checkFilters(input os.FileInfoEx) bool {
	for _, filter := range f.filters {
		matched, _ := filter.callback(input)
		if !matched {
			return false
		}
	}
	return true
}

func (f *Finder) addFilter(callback func(fiex os.FileInfoEx) (bool, error), order int) {
	f.filters = append(f.filters, struct {
		callback func(ex os.FileInfoEx) (bool, error)
		order  int
	}{callback, order})
	sort.Slice(f.filters, func(i, j int) bool {
		return f.filters[i].order < f.filters[j].order
	})
}

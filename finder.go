package finder

import (
	"github.com/pkg/errors"
	"github.com/duffpl/go-finder/mimechecker"
	"sort"
	"github.com/bmatcuk/doublestar"
	"sync"
	"github.com/duffpl/go-finder/checksum"
	"github.com/duffpl/go-finder/file"
)

type Finder struct {
	numCheckers int
	globFunc    FileInfoExGlobFunc
	filters []struct {
		callback func(ex file.FileInfoEx) (bool, error)
		order    int
	}
	lastErr error
}

var (
	defaultFileInfoExGlob            FileInfoExGlobFunc
	defaultFilterCheckersConcurrency int = 8
)

func init() {
	mc := mimechecker.NewMulti(mimechecker.NewGoHttp(), mimechecker.NewGoMime())
	defaultFileInfoExGlob = NewLazyGlobber(doublestar.Glob, checksum.MD5ByPath, mc.TypeByFile)
}

func New() *Finder {
	return new(Finder).
		SetGlobFunc(defaultFileInfoExGlob).
		SetCheckerConcurrency(defaultFilterCheckersConcurrency)
}

// SetGlobFunc sets function that is used for fetching list of FileInfoEx items used for filters. Default glob function
// is one that creates "lazy" FileInfoExs and uses doublestar.Glob (https://github.com/bmatcuk/doublestar) as lister since
// default Go globber doesn't repeat dir separator when using double asterisk
func (f *Finder) SetGlobFunc(gf FileInfoExGlobFunc) *Finder {
	f.globFunc = gf
	return f
}

// SetCheckerConcurrency sets number of go routines that are used for checking filters. Default is 8. Usually I/O will
// be bottleneck.
func (f *Finder) SetCheckerConcurrency(num int) *Finder {
	if num <= 0 {
		f.lastErr = errors.New("checkers count must be larger than 0")
	} else {
		f.numCheckers = num
	}
	return f
}

// Glob is used for fetching result list of files as slice of file.FileInfoEx items
func (f *Finder) Glob(pattern string) (result []file.FileInfoEx, err error) {
	if f.lastErr != nil {
		err = f.lastErr
		return
	}
	var globResult []file.FileInfoEx
	if globResult, err = f.globFunc(pattern); err != nil {
		err = errors.Wrap(err, "glob")
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(globResult))
	filtersOutput := f.runFilters(16, wg, globResult)
	createConsumer(filtersOutput, &result, wg)
	wg.Wait()
	return
}

func (f *Finder) createWorkers(cnt int, in <-chan file.FileInfoEx, out chan<- file.FileInfoEx, itemCount int, consumerWaitGroup *sync.WaitGroup) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(itemCount)
	for i:=0;i<cnt;i++ {
		go func() {
			for info := range in {
				if f.checkFilters(info) {
					consumerWaitGroup.Add(1)
					out <- info
				}
				wg.Done()
			}
		}()
	}
	return wg
}

func (f *Finder) runFilters(workerCnt int, tasksWg *sync.WaitGroup, entries []file.FileInfoEx) (chan file.FileInfoEx) {
	output := make(chan file.FileInfoEx)
	in := make(chan file.FileInfoEx)
	for i:=0;i<workerCnt;i++ {
		go func() {
			for info := range in {
				if f.checkFilters(info) {
					output <- info
				} else {
					tasksWg.Done()
				}
			}
		}()
	}
	go func() {
		for _, entry := range entries {
			in <- entry
		}
		close(in)
	}()
	return output
}


func createConsumer(in chan file.FileInfoEx, result *[]file.FileInfoEx, tasksWg *sync.WaitGroup) {
	go func() {
		for info := range in {
			*result = append(*result, info)
			tasksWg.Done()
		}
	}()
}

func (f *Finder) checkFilters(input file.FileInfoEx) bool {
	for _, filter := range f.filters {
		matched, _ := filter.callback(input)
		if !matched {
			return false
		}
	}
	return true
}

func (f *Finder) addFilter(callback func(fiex file.FileInfoEx) (bool, error), order int) {
	f.filters = append(f.filters, struct {
		callback func(ex file.FileInfoEx) (bool, error)
		order    int
	}{callback, order})
	sort.Slice(f.filters, func(i, j int) bool {
		return f.filters[i].order < f.filters[j].order
	})
}

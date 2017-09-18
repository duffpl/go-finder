package finder

import (
	"regexp"
	"fmt"

	"github.com/pkg/errors"
	"github.com/duffpl/go-finder/file"
)
// CmpOperator is string const enum for Size filter
type CmpOperator string

const (
	MoreThan CmpOperator = ">"
	MoreOrEqual = ">="
	LessThan = "<"
	LessOrEqual = "<="
	Equal = "=="
)

var Errors = struct{
	InvalidSizeOperator error
}{
	InvalidSizeOperator: errors.New("invalid size operator"),
}

// Checksum adds matching against checksum. Expected checksum should be hex encoded string
func (f *Finder) Checksum(hexChecksum string) *Finder {
	if f.lastErr != nil { return f }
	f.addFilter(func(fiex file.FileInfoEx) (result bool, err error) {
		var fileChecksum []byte
		if fileChecksum, err = fiex.Checksum(); err != nil {
			err = errors.Wrap(err, "checksum")
			return
		}
		result = fmt.Sprintf("%x", fileChecksum) == hexChecksum
		return
	}, 100)
	return f
}

func isCmpOperatorValid(cmpOp CmpOperator) bool {
	validOperators := []CmpOperator{MoreThan, MoreOrEqual, LessThan, LessOrEqual, Equal}
	for _, vop := range validOperators {
		if vop == cmpOp {
			return true
		}
	}
	return false
}

// Size adds matching against file size. First argument is operator type. Valid operators are available in
// CmpOperator const. Returns error if operator isn't allowed.
// Chains as AND operator
func (f *Finder) Size(cmpOp CmpOperator, cmpSize int64) *Finder {
	if f.lastErr != nil { return f }
	if !isCmpOperatorValid(cmpOp) {
		f.lastErr = Errors.InvalidSizeOperator
		return f
	}
	f.addFilter(func(info file.FileInfoEx) (cmpResult bool, err error) {
		size := info.Size()
		switch cmpOp {
		case MoreThan:
			cmpResult = size > cmpSize
		case MoreOrEqual:
			cmpResult = size >= cmpSize
		case LessThan:
			cmpResult = size < cmpSize
		case LessOrEqual:
			cmpResult = size <= cmpSize
		case Equal:
			cmpResult = size == cmpSize
		}
		return
	}, 1)
	return f
}

// Mime adds matching against MIME type of file
func (f *Finder) Mime(mimeType string) *Finder {
	if f.lastErr != nil { return f }
	f.addFilter(func(ex file.FileInfoEx) (result bool, err error) {
		var mimeResult string
		if mimeResult, err = ex.Mime(); err != nil {
			return
		}
		return mimeResult == mimeType, nil
	}, 50)
	return f
}

// MimeRegexp adds matching against MIME type of file using regexp pattern. This is useful for finding files of given
// type. E.g. to find all images > MimeRegexp("^image")
func (f *Finder) MimeRegexp(pattern string) *Finder {
	if f.lastErr != nil { return f }
	var compiled *regexp.Regexp
	if compiled, f.lastErr = regexp.Compile(pattern); f.lastErr != nil {
		return f
	}
	f.addFilter(func(ex file.FileInfoEx) (result bool, err error) {
		var mimeResult string
		if mimeResult, err = ex.Mime(); err != nil {
			return
		}
		return compiled.Match([]byte(mimeResult)), nil
	}, 50)
	return f
}

// RegexpName adds matching against file name using regexp pattern.
func (f *Finder) RegexpName(pattern string) *Finder {
	if f.lastErr != nil { return f }
	var compiled *regexp.Regexp
	if compiled, f.lastErr = regexp.Compile(pattern); f.lastErr != nil {
		return f
	}
	f.addFilter(func(ex file.FileInfoEx) (bool, error) {
		return compiled.Match([]byte(ex.Name())), nil
	}, 2)
	return f
}

// RegexpName adds matching against full file path using regexp pattern.
func (f *Finder) RegexpPath(pattern string) *Finder {
	if f.lastErr != nil { return f }
	var compiled *regexp.Regexp
	if compiled, f.lastErr = regexp.Compile(pattern); f.lastErr != nil {
		return f
	}
	f.addFilter(func(ex file.FileInfoEx) (result bool, err error) {
		var abs string
		if abs, err = ex.Abs(); err != nil {
			err = errors.Wrap(err, "RegexpPath")
		} else {
			result = compiled.Match([]byte(abs))
		}
		return
	}, 2)
	return f
}
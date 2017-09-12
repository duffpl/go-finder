package finder

import (
	"github.com/duffpl/go-rgxp"
	"github.com/pkg/errors"
	"regexp"
	"github.com/duffpl/go-finder/os"
)

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

func isCmpOperatorValid(cmpOp CmpOperator) bool {
	validOperators := []CmpOperator{MoreThan, MoreOrEqual, LessThan, LessOrEqual, Equal}
	for _, vop := range validOperators {
		if vop == cmpOp {
			return true
		}
	}
	return false
}

func (f *Finder) Size(cmpOp CmpOperator, cmpSize int64) *Finder {
	if f.err != nil {
		return f
	}
	if !isCmpOperatorValid(cmpOp) {
		f.err = Errors.InvalidSizeOperator
		return f
	}
	f.addFilter(func(info os.FileInfoEx) (cmpResult bool, err error) {
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

func (f *Finder) Mime(mimeType string) *Finder {
	f.addFilter(func(ex os.FileInfoEx) (result bool, err error) {
		var mimeResult string
		if mimeResult, err = ex.Mime(); err != nil {
			return
		}
		return mimeResult == mimeType, nil
	}, 50)
	return f
}

func (f *Finder) MimeRegexp(pattern string) *Finder {
	if f.err != nil { return f }
	var compiled *regexp.Regexp
	if compiled, f.err = regexp.Compile(pattern); f.err != nil {
		return f
	}
	f.addFilter(func(ex os.FileInfoEx) (result bool, err error) {
		var mimeResult string
		if mimeResult, err = ex.Mime(); err != nil {
			return
		}
		return compiled.Match([]byte(mimeResult)), nil
	}, 50)
	return f
}

func (f *Finder) RegexpName(patterns []string) *Finder {
	var compiled []*regexp.Regexp
	if compiled, f.err = rgxp.CompileAll(patterns); f.err != nil {
		return f
	}
	f.addFilter(func(ex os.FileInfoEx) (bool, error) {
		return rgxp.MatchAny(compiled, ex.Name()), nil
	}, 2)
	return f
}

func (f *Finder) RegexpPath(patterns []string) *Finder {
	var compiled []*regexp.Regexp
	if compiled, f.err = rgxp.CompileAll(patterns); f.err != nil {
		return f
	}
	f.addFilter(func(ex os.FileInfoEx) (bool, error) {
		if abs, err := ex.Abs(); err != nil {
			return false, errors.Wrap(err, "RegexpPath")
		} else {
			return rgxp.MatchAny(compiled, abs), nil
		}
	}, 2)
	return f
}


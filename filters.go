package finder

import (
	"github.com/duffpl/go-rgxp"
	"github.com/pkg/errors"
	"path/filepath"
)

type CmpOperator string

const (
	MoreThan CmpOperator = ">"
	MoreOrEqual = ">="
	LessThan = "<"
	LessOrEqual = "<="
	Equal = "=="
)

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
		f.err = errors.New("invalid size compare operator")
		return f
	}
	f.addFilter(func(info *FileInfoEx) (cmpResult bool, err error) {
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
	})
	return f
}

func (f *Finder) Mime(mimeType string) *Finder {
	f.addFilter(func(ex *FileInfoEx) (bool, error) {
		mimeResult, err := stdMimeChecker.ByPath(ex.AbsolutePath)
		if err != nil {
			return false, err
		}
		return mimeResult == mimeType, nil
	})
	return f
}

func (f *Finder) Match(pattern string) *Finder {
	f.addFilter(func(ex *FileInfoEx) (bool, error) {
		return filepath.Match(pattern, ex.Name())
	})
	return f
}

func (f *Finder) RegexpName(patterns []string) *Finder {
	compiled, err := rgxp.CompileAll(patterns)
	if err != nil {
		f.err = err
		return f
	}
	f.addFilter(func(ex *FileInfoEx) (bool, error) {
		return rgxp.MatchAny(compiled, ex.Name()), nil
	})
	return f
}

func (f *Finder) RegexpPath(patterns []string) *Finder {
	compiled, err := rgxp.CompileAll(patterns)
	if err != nil {
		f.err = err
		return f
	}
	f.addFilter(func(ex *FileInfoEx) (bool, error) {
		return rgxp.MatchAny(compiled, ex.AbsolutePath), nil
	})
	return f
}


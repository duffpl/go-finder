package finder

import (
	"github.com/bmatcuk/doublestar"
)

type globber interface {
	Glob(pattern string) (matches []string, err error)
}

type doubleStarGlobber struct {}

func (doubleStarGlobber) Glob(pattern string) (matches []string, err error) {
	return doublestar.Glob(pattern)
}

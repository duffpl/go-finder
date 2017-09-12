package mimechecker

import (
	"mime"
	"path/filepath"
)

type GoMime struct {}

func (GoMime) TypeByFile(path string) (string, error) {
	return mime.TypeByExtension(filepath.Ext(path)), nil
}

func NewGoMime() *GoMime {
	return &GoMime{}
}
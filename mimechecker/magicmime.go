package mimechecker

import (
	"github.com/rakyll/magicmime"
)

type magicMime struct {}

func (*magicMime) ByPath(path string) (m string, err error) {
	return magicmime.TypeByFile(path)
}

func NewMagicMime() *magicMime {
	return &magicMime{}
}
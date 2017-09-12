package os

import (
	"os"
)

type FileInfoEx interface {
	os.FileInfo
	Abs() (abs string, err error)
	Checksum() (cs []byte, err error)
	Mime() (m string, err error)
}


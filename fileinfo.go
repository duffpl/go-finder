package finder

import (
	"os"
)

type FileInfoEx struct {
	os.FileInfo
	AbsolutePath string
	RelativePath string
	Checksum []byte
}



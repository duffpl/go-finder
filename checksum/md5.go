package checksum

import (
	"crypto/md5"
	"io"
	"os"
	"github.com/pkg/errors"
)

func MD5ByPath(path string) (checksum []byte, err error) {
	var handle *os.File
	checksumWriter := md5.New()
	if handle, err = os.Open(path); err != nil {
		err = errors.Wrap(err, "os.Open")
		return
	}
	if _, err = io.Copy(checksumWriter, handle); err != nil {
		err = errors.Wrap(err, "io.Copy")
		return
	}
	defer handle.Close()
	checksum = checksumWriter.Sum(nil)
	return
}

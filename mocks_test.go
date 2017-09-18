package finder

import (
	"time"
	"os"
	"github.com/duffpl/go-finder/file"
)

type mockFileInfoEx struct {
	name     string
	size     int64
	mime     string
	abs      string
	checksum []byte
}

func (m *mockFileInfoEx) Name() string {
	return m.name
}

func (m *mockFileInfoEx) Size() int64 {
	return m.size
}

func (*mockFileInfoEx) Mode() os.FileMode {
	return 0
}

func (*mockFileInfoEx) ModTime() time.Time {
	return time.Now()
}

func (*mockFileInfoEx) IsDir() bool {
	return false
}

func (*mockFileInfoEx) Sys() interface{} {
	return struct{}{}
}

func (m *mockFileInfoEx) Abs() (abs string, err error) {
	return m.abs, nil
}

func (m *mockFileInfoEx) Checksum() (cs []byte, err error) {
	return m.checksum, nil
}

func (m *mockFileInfoEx) Mime() (r string, err error) {
	return m.mime, nil
}

func newMockGlobFunc(result []file.FileInfoEx) FileInfoExGlobFunc {
	return func(pattern string) ([]file.FileInfoEx, error) {
		return result, nil
	}
}
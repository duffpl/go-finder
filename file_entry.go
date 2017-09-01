package finder

import (
	"os"
	"time"
	"fmt"
	"crypto/md5"
)

type FileEntry struct {
	Mode os.FileMode
	IsDir bool
	Name string
	Size int64
	ModTime time.Time
	AbsolutePath string
	RelativePath string
	Checksum []byte
	_hash []byte
}

func (f *FileEntry) Hash() []byte {
	if f._hash == nil {
		md := md5.New()
		f._hash = md.Sum([]byte(fmt.Sprintf("%s%s",f.RelativePath, f.Checksum)))
	}
	return f._hash
}

func newFileEntryFromFileInfo(fi os.FileInfo) *FileEntry {
	return &FileEntry{
		IsDir: fi.IsDir(),
		ModTime: fi.ModTime(),
		Size: fi.Size(),
		Name: fi.Name(),
		Mode: fi.Mode(),
	}
}


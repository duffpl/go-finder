package mimechecker

import (
	"os"
	"errors"
	"net/http"
)

type goHttp struct {}

const mimeOctet = "application/octet-stream"

func (*goHttp) ByPath(path string) (m string, err error) {
	defer func() {
		if err != nil {
			err = errors.New("cannot read mime from file:" + err.Error())
		}
	}()
	buf := make([]byte, 512)
	fh, err := os.Open(path)
	if err != nil {
		return
	}
	_, err = fh.Read(buf)
	if err != nil {
		return
	}
	m = http.DetectContentType(buf)
	if m == mimeOctet {
		m = ""
	}
	return
}

func NewGoHttp() *goHttp {
	return &goHttp{}
}
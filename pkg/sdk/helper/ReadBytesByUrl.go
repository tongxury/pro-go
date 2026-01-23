package helper

import (
	"io"
	"net/http"
)

func ReadBytesByUrl(url string) (bytes []byte, err error) {

	get, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(get.Body)
}

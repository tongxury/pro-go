package imagez

import (
	"io"
	"net/http"
)

func GetBlob(url string) ([]byte, error) {

	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	return all, nil

}

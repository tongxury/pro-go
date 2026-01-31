package bytez

import (
	"fmt"
	"io"
	"net/http"
)

func ReadFileBytes(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to load url: %s, status: %d", url, resp.StatusCode)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(imgBytes) == 0 {
		return nil, fmt.Errorf("url (%s) is empty", url)
	}

	return imgBytes, nil
}

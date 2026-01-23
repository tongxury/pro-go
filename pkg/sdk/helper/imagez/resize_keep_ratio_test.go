package imagez

import (
	"os"
	"testing"
)

func TestResizeKeepRatio(t *testing.T) {

	file, err := os.ReadFile("01b7a8b6c6bbdc1583bae7b5f8dad6ca.jpeg")
	if err != nil {
		return
	}

	ratio, err := ResizeKeepRatio(file, 720, 1280)
	if err != nil {
		return
	}

	create, err := os.Create("tmp.jpeg")
	if err != nil {
		return
	}

	create.Write(ratio)
}

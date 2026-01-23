package imagez

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func Resize(srcBytes []byte, width, height int) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(srcBytes))
	if err != nil {
		return nil, err
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 90})
	case "png":
		err = imaging.Encode(&buf, resized, imaging.PNG)
	default:
		// 也可以统一用 imaging.Encode 并指定想要的格式
		err = imaging.Encode(&buf, resized, imaging.JPEG)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

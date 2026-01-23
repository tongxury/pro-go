package imagez

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
)

// Split3x3 splits an image bytes into a 3x3 grid (9 equal parts).
// It returns a slice of 9 images (top-left, top-center, top-right, mid-left, ...).
func Split3x3(data []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("image data is empty")
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < 3 || height < 3 {
		return nil, errors.New("image is too small to split into 3x3")
	}

	// Calculate cell size
	cellWidth := width / 3
	cellHeight := height / 3

	var result [][]byte

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {

			// Calculate crop rectangle
			x0 := bounds.Min.X + col*cellWidth
			y0 := bounds.Min.Y + row*cellHeight
			x1 := x0 + cellWidth
			y1 := y0 + cellHeight

			// Fix resolution loss: Ensure the last column/row covers the remainder pixels
			if col == 2 {
				x1 = bounds.Max.X
			}
			if row == 2 {
				y1 = bounds.Max.Y
			}

			rect := image.Rect(x0, y0, x1, y1)

			var sub image.Image

			// Use SubImage if supported for performance (shares buffer)
			if interfaceImg, ok := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}); ok {
				sub = interfaceImg.SubImage(rect)
			} else {
				// Fallback: create new image and draw
				dst := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
				draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
				sub = dst
			}

			var buf bytes.Buffer
			var err error
			switch format {
			case "jpeg":
				err = jpeg.Encode(&buf, sub, &jpeg.Options{Quality: 100})
			case "png":
				err = png.Encode(&buf, sub)
			default:
				// Fallback to PNG if unknown format or not explicitly handled
				err = png.Encode(&buf, sub)
			}

			if err != nil {
				return nil, err
			}
			result = append(result, buf.Bytes())
		}
	}

	return result, nil
}

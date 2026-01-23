package imagez

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

func ResizeKeepRatio(src []byte, w, h int) ([]byte, error) {
	// 解码原始图片
	img, format, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	// 计算原始尺寸
	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// 计算保持宽高比的缩放尺寸
	ratio := float64(srcW) / float64(srcH)
	targetRatio := float64(w) / float64(h)

	var scaleW, scaleH int
	var offsetX, offsetY int

	if ratio > targetRatio {
		// 原图更宽，以高度为准缩放
		scaleH = h
		scaleW = int(float64(h) * ratio)
		offsetX = (scaleW - w) / 2
		offsetY = 0
	} else {
		// 原图更高，以宽度为准缩放
		scaleW = w
		scaleH = int(float64(w) / ratio)
		offsetX = 0
		offsetY = (scaleH - h) / 2
	}

	// 先缩放到中间尺寸
	scaled := image.NewRGBA(image.Rect(0, 0, scaleW, scaleH))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), img, srcBounds, draw.Over, nil)

	// 创建最终的目标图片（精确尺寸 w x h）
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	// 从缩放后的图片中裁剪出中心部分
	srcRect := image.Rect(offsetX, offsetY, offsetX+w, offsetY+h)
	draw.Draw(dst, dst.Bounds(), scaled, srcRect.Min, draw.Over)

	// 编码为字节数组
	var buf bytes.Buffer
	switch format {
	case "png":
		err = png.Encode(&buf, dst)
	case "jpeg":
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 95})
	default:
		// 默认使用 JPEG
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 95})
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

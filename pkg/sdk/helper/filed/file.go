package filed

import "strings"

func FindSuffix(name string) string {
	// 转换为小写以便不区分大小写比较
	lowerName := strings.ToLower(name)

	// 视频格式
	videoFormats := []string{
		"mp4", "avi", "mov", "wmv", "flv", "mkv", "webm", "m4v", "mpeg", "mpg", "3gp", "ts", "rm", "rmvb", "vob", "asf",
	}

	// 图片格式
	imageFormats := []string{
		"jpg", "jpeg", "png", "gif", "bmp", "webp", "svg", "tiff", "tif", "heic", "heif", "raw", "cr2", "nef", "arw", "dng",
	}

	// 检查视频格式
	for _, format := range videoFormats {
		if strings.HasSuffix(lowerName, format) {
			return "." + format
		}
	}

	// 检查图片格式
	for _, format := range imageFormats {
		if strings.HasSuffix(lowerName, format) {
			return "." + format
		}
	}

	// 如果没有匹配到任何格式，返回空字符串
	return ""
}

func ContentTypeByName(name string) string {

	t := "application/octet-stream"

	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		tp := parts[len(parts)-1]

		// 图片

		if strings.Contains(tp, "png") {
			return "image/png"
		}

		if strings.Contains(tp, "jpg") {
			return "image/jpeg"
		}

		if strings.Contains(tp, "jpeg") {
			return "image/jpeg"
		}

		if strings.Contains(tp, "gif") {
			return "image/gif"
		}

		if strings.Contains(tp, "bmp") {
			return "image/bmp"
		}

		// 文件
		if strings.Contains(tp, "pdf") {
			return "application/pdf"
		}

		if strings.Contains(tp, "docx") {
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		}

		if strings.Contains(tp, "doc") {
			return "application/msword"
		}
	}

	return t
}

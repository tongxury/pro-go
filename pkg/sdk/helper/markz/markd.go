package markz

import "strings"

func Verify(src string) string {

	src = strings.ReplaceAll(src, ".", "\\.")
	src = strings.ReplaceAll(src, "{", "\\{")
	src = strings.ReplaceAll(src, "}", "\\}")
	src = strings.ReplaceAll(src, "-", "\\-")
	src = strings.ReplaceAll(src, "_", "\\_")
	src = strings.ReplaceAll(src, "+", "\\+")
	src = strings.ReplaceAll(src, "(", "\\(")
	src = strings.ReplaceAll(src, ")", "\\)")
	src = strings.ReplaceAll(src, "$", "\\$")
	src = strings.ReplaceAll(src, "=", "\\=")

	return src
}

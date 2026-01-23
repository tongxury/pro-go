package phonenum

import (
	"regexp"
)

func IsValid(number string) bool {
	//return regexp.MustCompile("^1[345789]{1}\\d{9}$").MatchString(number)
	return regexp.MustCompile("^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$").MatchString(number)
}

package redizv2

import (
	"fmt"
)

func String(val any) string {
	return fmt.Sprintf("\"%v\"", val)
}

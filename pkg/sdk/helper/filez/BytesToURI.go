package filez

import (
	"encoding/base64"
	"fmt"
)

func BytesToURI(data []byte) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:image/png;base64,%s", b64)
}

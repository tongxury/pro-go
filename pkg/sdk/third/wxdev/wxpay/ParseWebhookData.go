package wxpay

import (
	"context"
	"encoding/base64"
	"fmt"
)

func (t *Client) ParseWebhookData(ctx context.Context) string {

	decodeString, err := base64.StdEncoding.DecodeString("0CPAOyQbldtz16Drdx8/wZUH4HWk/9enTYkMcCqyk3vzqewCEZMWu3JCWIIWeiCCU30GeX+dhVku6o5FmUemHLGqbCzB3DeClbWmon3k0s8VquPi3DTHLOOHi5w4nlptv61hCjAY4msFL8q376CytSH1bMUq/5yPF7HlUyegJLpjcntT/vms1sVwbRX0YhlZPCmevSUGqOnSTZcSEG0opu06Q4F9JvdfFXGBqglWI74+D51rokIE92RVa79lFzaqHPXYCEFJAjFPbZQGYD4jWhzvYeBVFU22k5eiie/4uUjzj/rAgy58QEIFH3l8piqioaYhh/8rpRASCjGvpEmr0clz/AHYIb30axkt0MZx44qbAs3BcvPDYeGAdnrBzIYOWFrOiwTyhqh4oKxkrDFo65ZREEiPe1kdLRAtRD/4anc1wIYw0alubC0oam1Q69SmVXKqh818/8GMYXgnyj6fjYhSRvkkRwNeRpcdvGbxMj8p2XxhI1QWpcpVR6zuHi74K3kapiE71TZBqmQjDeaSQ+bOKBspn++fYaYOTBj1VO214SjJwgFnVQwzYaioqOxep1qSCq1BGQ==")
	if err != nil {
		return ""
	}

	fmt.Println(string(decodeString))

	return ""
}

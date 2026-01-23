package dypenapi

import (
	"context"
	"store/confs"
	"testing"
)

func TestClient_GetAuthCode(t *testing.T) {

	c, _ := NewClient(Config{
		ClientKey:    confs.DouyinClientKey,
		ClientSecret: confs.DouyinClientSecret,
	})

	//c.getAccessToken()
	c.GetAuthCodeQrcode(context.Background())
}

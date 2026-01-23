package vengine

import (
	"store/confs"

	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

// https://github.com/volcengine/volcengine-go-sdk/blob/master/SDK_Integration_zh.md
type Client struct {
	sess *session.Session
	conf Config
}

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
}

const (
	// 请求凭证，从访问控制申请
	AccessKeyID     = confs.BytedanceAccessKeyID
	SecretAccessKey = confs.BytedanceSecretAccessKey

	// 请求地址
	Addr = "https://icp.volcengineapi.com"
	Path = "/" // 路径，不包含 Query

	// 请求接口信息
	Service = "iccloud_muse"
	Region  = "cn-north"
	Action  = "SearchTemplate"
	Version = "2021-09-01"
)

func NewClient() (*Client, error) {

	conf := Config{
		AccessKeyID:     confs.BytedanceAccessKeyID,
		AccessKeySecret: confs.BytedanceSecretAccessKey,
	}

	region := "cn-north"
	config := volcengine.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(conf.AccessKeyID, conf.AccessKeySecret, "")).
		WithRegion(region)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		sess: sess,
		conf: conf,
	}, nil
}

func (t *Client) C() error {

	return nil
}

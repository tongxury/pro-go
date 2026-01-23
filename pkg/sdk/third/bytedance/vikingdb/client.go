package vikingdb

import "store/confs"

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	BaseURL         string
	Service         string
	Region          string
	Host            string
}
type Client struct {
	conf Config
}

func NewClient() *Client {

	return &Client{
		conf: Config{
			AccessKeyID:     confs.BytedanceAccessKeyID,
			AccessKeySecret: confs.BytedanceSecretAccessKey,

			BaseURL: "https://api-vikingdb.vikingdb.cn-beijing.volces.com",

			Service: "vikingdb",
			Region:  "cn-beijing",
			Host:    "api-vikingdb.vikingdb.cn-beijing.volces.com",
		},
	}
}

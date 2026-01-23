package confcenter

import (
	"store/pkg/clients/ch"
	"store/pkg/clients/elastics"
	"store/pkg/clients/mgz"
	"store/pkg/clients/mysqlz"
	"store/pkg/clients/redizv2"
	"store/pkg/rediz"
	"store/pkg/sdk/third/aliyun/alioss"
	"store/pkg/sdk/third/aws/s3"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/qiniu"
	"time"
)

// config 结构体必须定义在confcenter包中，否则无法generate成功
type Mysql struct {
	Driver string
	Source string
}

type Redis struct {
	// host:port address.
	Addrs        []string
	Password     string
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type Elasticsearch struct {
	Addresses []string
	Username  string
	Password  string
}

type Mongo struct {
	Uri           string
	AuthMechanism string
	Username      string
	Password      string
	Database      string
}

/*
database:

	clickhouse:
	  addrs:
	    - xx.xx.xx.xx:9000
	  database: default
	  username: default
	  password:
*/
type ClickHouseConfig struct {
	Addrs    []string
	Database string
	Username string
	Password string
}

type Database struct {
	Mysql        Mysql
	Mysqlz       mysqlz.Config
	RedisCluster redizv2.ClusterConfig
	//RedisV2       redizv2.Config
	Redis         Redis
	Rediz         rediz.Config
	Mongo         mgz.Config
	Clickhouse    ClickHouseConfig
	Elasticsearch Elasticsearch
	S3            s3.Config
	Oss           alioss.Config
	Tos           tos.Config
	Qiniu         qiniu.Config
	Ch            ch.Config
	Elastics      elastics.Config
}

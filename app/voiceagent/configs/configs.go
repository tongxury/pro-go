package configs

import (
	"fmt"
	"os"
	"store/confs"
	"store/pkg/clients/grpcz"

	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/aws/s3"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/qiniu"
	"time"

	"entgo.io/ent/dialect"
	"github.com/go-kratos/kratos/v2/log"
)

var ip = "101.132.192.41"

type BizConfig struct {
}

var configs = map[string]*confcenter.Config[BizConfig]{
	"prod_qiniu": {
		Meta: confcenter.Meta{
			Namespace: "prod",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8080",
				Timeout: 300 * time.Second,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8090",
				Timeout: 300 * time.Second,
			},
		},
		Logger: log.DefaultLogger,
		Database: confcenter.Database{
			Mongo: mgz.Config{
				Uri:      fmt.Sprintf("mongodb://%s:27017/pro?retryWrites=true&w=majority", "101.132.192.41"),
				Username: "root",
				Password: confs.MongoPassword,
				Database: "pro",
			},
			Rediz: rediz.Config{
				Addrs:    []string{"101.132.192.41:6379"},
				Password: confs.RedisPassword,
			},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,
				Endpoint:        "tos-cn-shanghai.volces.com",
				DefaultBucket:   "veres",
				Region:          "cn-shanghai",
			},
		},
		Component: confcenter.Component{
			Grpc: grpcz.Configs{
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter.prod.svc.cluster.local:8090",
				},
			},
			Genai: gemini.FactoryConfig{
				Configs: []*gemini.Config{
					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.GeminiKeys[23]},
					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.GeminiKeys[18]},
				},
			},
		},
		Biz: BizConfig{},
	},
	"prod-cn": {
		Meta: confcenter.Meta{
			Namespace: "prod",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8080",
				Timeout: 300 * time.Second,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8090",
				Timeout: 300 * time.Second,
			},
		},
		Logger: log.DefaultLogger,
		Database: confcenter.Database{
			Mongo: mgz.Config{
				//Uri:      "mongodb://mongodb-sharded.prod:27017/tgbot?retryWrites=true&w=majority",
				//Username: "root",
				//Password: "z4XNmlaOjo",
				Uri:      fmt.Sprintf("mongodb://%s:27017/pro?retryWrites=true&w=majority", "101.132.192.41"),
				Username: "root",
				Password: confs.MongoPassword,
				Database: "pro",
			},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"101.132.192.41:6379"},
				Password: confs.RedisPassword,
			},
			//Oss: alioss.Config{
			//	AccessSecret: "IGfLHuBHcTHYvYvEj3ZekmyKBG3wIu",
			//	Bucket:       "veogocn",
			//	Endpoint:     "oss-cn-beijing.aliyuncs.com",
			//	//Endpoint: "oscar-res-491f145810273c7d8c143a17a9ce5ed21a-ossalias",
			//},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,

				//Endpoint:        "tos-s3-cn-shanghai.volces.com",
				Endpoint: "tos-cn-shanghai.volces.com",
				//Endpoint: "tos-accelerate.volces.com",
				//yoozyres.tos-accelerate.volces.com
				DefaultBucket: "veres",
				Region:        "cn-shanghai",
			},

			S3: s3.Config{
				AccessKey:    confs.AWSS3AccessKey,
				AccessSecret: confs.AWSS3Secret,

				Region:   "cn-north-1",
				Bucket:   "veogoresources",
				Endpoint: "https://veogoresources.s3.cn-north-1.amazonaws.com.cn",
			},
			Qiniu: qiniu.Config{
				AccessKey:    confs.QiniuAccessKey,
				AccessSecret: confs.QiniuAccessSecret,
			},
		},
		Component: confcenter.Component{
			Grpc: grpcz.Configs{
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter.prod.svc.cluster.local:8090",
				},
			},
			Genai: gemini.FactoryConfig{
				Configs: []*gemini.Config{
					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.GeminiKeys[23]},
					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.GeminiKeys[18]},
				},
			},
		},
		Biz: BizConfig{},
	},
	"prod": {
		Meta: confcenter.Meta{
			Namespace: "prod",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8080",
				Timeout: 300 * time.Second,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8090",
				Timeout: 300 * time.Second,
			},
		},
		Logger: log.DefaultLogger,
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordProj, ip),
			},
			Clickhouse: confcenter.ClickHouseConfig{
				Addrs:    []string{ip + ":9000"},
				Database: "default",
				Username: "default",
				Password: confs.ClickhousePassword,
			},
			Mongo: mgz.Config{
				//Uri:      "mongodb://mongodb-sharded.prod:27017/tgbot?retryWrites=true&w=majority",
				//Username: "root",
				//Password: "z4XNmlaOjo",
				Uri:      fmt.Sprintf("mongodb://%s:27017/pro?retryWrites=true&w=majority", "101.132.192.41"),
				Username: "root",
				Password: confs.MongoPassword,

				Database: "pro",
			},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"101.132.192.41:6379"},
				Password: confs.RedisPassword,
			},
			//Oss: alioss.Config{
			//	AccessSecret: "IGfLHuBHcTHYvYvEj3ZekmyKBG3wIu",
			//	Bucket:       "veogocn",
			//	Endpoint:     "oss-cn-beijing.aliyuncs.com",
			//	//Endpoint: "oscar-res-491f145810273c7d8c143a17a9ce5ed21a-ossalias",
			//},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,

				//Endpoint:        "tos-s3-cn-shanghai.volces.com",
				Endpoint: "tos-cn-shanghai.volces.com",
				//Endpoint: "tos-accelerate.volces.com",
				//yoozyres.tos-accelerate.volces.com
				DefaultBucket: "veres",
				Region:        "cn-shanghai",
			},

			Qiniu: qiniu.Config{
				AccessKey:    confs.QiniuAccessKey,
				AccessSecret: confs.QiniuAccessSecret,
			},
		},
		Component: confcenter.Component{
			Kafka: confcenter.KafkaConfig{
				Brokers: []string{"kafka-headless.prod.svc.cluster.local:9092"},
			},
			Solana: confcenter.SolanaConfig{
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			Grpc: grpcz.Configs{
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter.prod.svc.cluster.local:8090",
				},
			},
			Genai: gemini.FactoryConfig{
				Configs: []*gemini.Config{
					{Proxy: "", Key: confs.GeminiKeys[23]},
					{Proxy: "", Key: confs.GeminiKeys[18]},
				},
			},
		},
		Biz: BizConfig{},
	},
	"dev": {
		Meta: confcenter.Meta{
			Namespace: "dev",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8086",
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8096",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mongo: mgz.Config{
				//Uri:      "mongodb://mongodb-sharded.prod:27017/tgbot?retryWrites=true&w=majority",
				//Username: "root",
				//Password: "z4XNmlaOjo",
				Uri:      fmt.Sprintf("mongodb://%s:27017/pro?retryWrites=true&w=majority", "101.132.192.41"),
				Username: "root",
				Password: confs.MongoPassword,

				Database: "pro",
			},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"101.132.192.41:6379"},
				Password: confs.RedisPassword,
			},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,

				Endpoint:      "tos-cn-shanghai.volces.com",
				DefaultBucket: "veres",
				Region:        "cn-shanghai",
			},
			Qiniu: qiniu.Config{
				AccessKey:    confs.QiniuAccessKey,
				AccessSecret: confs.QiniuAccessSecret,
			},
		},
		Component: confcenter.Component{
			Solana: confcenter.SolanaConfig{
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			Grpc: grpcz.Configs{
				User: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
			},
		},
		Biz: BizConfig{},
	},
}

func GetConfig() (*confcenter.Config[BizConfig], error) {
	env := os.Getenv("POD_NAMESPACE")
	if env == "" {
		env = "dev"
	}

	cc, ok := configs[env]
	if !ok {
		return configs["dev"], nil
	}

	return cc, nil
}

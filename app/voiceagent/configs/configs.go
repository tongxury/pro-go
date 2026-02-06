package configs

import (
	"os"
	"store/confs"
	"store/pkg/clients/grpcz"

	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/qiniu"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var ip = "101.132.192.41"

type BizConfig struct {
}

var configs = map[string]*confcenter.Config[BizConfig]{
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
			Mongo: confs.MongoMy,
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"103.30.78.254:6379"},
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
			Kafka: confs.KafkaMy,
			Grpc: grpcz.Configs{
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter.prod.svc.cluster.local:8090",
				},
				VoiceAgent: &grpcz.Config{
					Endpoint: "voiceagent.prod.svc.cluster.local:8090",
				},
			},
			Genai: gemini.FactoryConfig{
				Configs: []*gemini.Config{
					{
						Project:         "yuzhi-483807", // 从 secrets.go 或 configs.go 中提取的 Project ID
						Location:        "us-central1",  // 常用 Vertex AI Location
						APIVersion:      "v1",
						CredentialsJSON: confs.VertexAiSecret,                             // 使用 JSON 字符串凭据
						Proxy:           "http://proxy:strOngPAssWOrd@45.78.194.147:6060", // 如有需要
					},
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
			Mongo: confs.MongoMy,
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
			Kafka: confcenter.KafkaConfig{
				Brokers: []string{ip + ":9092"},
			},
			Solana: confcenter.SolanaConfig{
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			Genai: gemini.FactoryConfig{
				Configs: []*gemini.Config{
					{
						Project:         "yuzhi-483807", // 从 secrets.go 或 configs.go 中提取的 Project ID
						Location:        "us-central1",  // 常用 Vertex AI Location
						APIVersion:      "v1",
						CredentialsJSON: confs.VertexAiSecret,                             // 使用 JSON 字符串凭据
						Proxy:           "http://proxy:strOngPAssWOrd@45.78.194.147:6060", // 如有需要
					},
				},
			},
			Grpc: grpcz.Configs{
				User: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
				VoiceAgent: &grpcz.Config{
					Endpoint: "localhost:8096",
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

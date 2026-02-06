package configs

import (
	"fmt"
	"os"
	"store/pkg/clients/elastics"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/aliyun/alioss"
	"store/pkg/sdk/third/aliyun/alisms"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/openaiz"

	"store/confs"

	"github.com/go-kratos/kratos/v2/log"
)

var serverConf = confcenter.Server{
	Http: &confcenter.ServerConfig{
		Addr:    "0.0.0.0:8089",
		Timeout: 30000000000,
	},
	Grpc: &confcenter.ServerConfig{
		Addr:    "0.0.0.0:8099",
		Timeout: 30000000000,
	},
}

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
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8090",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mongo: mgz.Config{
				//Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "yoozy_pro",
				//Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "admin",

				Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.MongoPasswordYuzhi,
				Database: "admin",
			},
			Elastics: elastics.Config{
				//Addresses: []string{"http://localhost:9200/"},
				//Addresses: []string{"http://118.196.63.209:9200/"},
				Addresses: []string{"http://yoozy-qv7.public.cn-hangzhou.es-serverless.aliyuncs.com:9200/"},
				Username:  "yoozy-qv7",
				Password:  confs.ElasticPassword,
			},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"118.196.63.209:6379"},
				Password: confs.RedisPasswordProj,
			},
			Oss: alioss.Config{
				AccessKey:    confs.AliyunOssAccessKey,
				AccessSecret: confs.AliyunOssSecret,

				Bucket:   "yoozy",
				Region:   "oss-cn-hangzhou",
				Endpoint: "oss-cn-hangzhou.aliyuncs.com",
			},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,

				//Endpoint:        "tos-s3-cn-shanghai.volces.com",
				Endpoint: "tos-cn-shanghai.volces.com",
				//Endpoint: "tos-accelerate.volces.com",
				//yoozyres.tos-accelerate.volces.com
				DefaultBucket: "yoozyres",
				Region:        "cn-shanghai",
			},
		},
		Component: confcenter.Component{
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
			OpenAI: openaiz.Config{
				AppKey: confs.OpenAIKeys[0],

				//Proxy:   "http://proxy:strOngPAssWOrd@45.78.194.147:6060",
				BaseUrl: "http://45.78.194.147:6000",
			},
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,

				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				ProjAdmin: &grpcz.Config{
					Endpoint: "admin:8090",
				},
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
				Credit: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
			},
		},
		Logger: log.DefaultLogger,
	},
	"beta": {
		Meta: confcenter.Meta{
			Namespace: "beta",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8080",
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8090",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mongo: mgz.Config{
				Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.MongoPasswordProj,

				Database: "yoozy_pro_beta",

				//Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "admin",

				//Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "admin",
			},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"118.196.63.209:6379"},
				Password: confs.RedisPasswordProj,

				DB: 1,
			},
			Oss: alioss.Config{
				AccessKey:    confs.AliyunOssAccessKey,
				AccessSecret: confs.AliyunOssSecret,

				Bucket:   "yoozy",
				Region:   "oss-cn-hangzhou",
				Endpoint: "oss-cn-hangzhou.aliyuncs.com",
			},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,

				//Endpoint:        "tos-s3-cn-shanghai.volces.com",
				Endpoint: "tos-cn-shanghai.volces.com",
				//Endpoint: "tos-accelerate.volces.com",
				//yoozyres.tos-accelerate.volces.com
				DefaultBucket: "yoozyres",
				Region:        "cn-shanghai",
			},
		},
		Component: confcenter.Component{
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
			OpenAI: openaiz.Config{
				AppKey: confs.OpenAIKeys[0],

				//Proxy:   "http://proxy:strOngPAssWOrd@45.78.194.147:6060",
				BaseUrl: "http://45.78.194.147:6000",
			},
			GetGoAPI: openaiz.Config{
				AppKey: confs.GetGoAPIKey,

				BaseUrl: "https://cn.getgoapi.com",
			},
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,

				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				ProjAdmin: &grpcz.Config{
					Endpoint: "admin:8090",
				},
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
				Credit: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
			},
		},
		Logger: log.DefaultLogger,
	},
	"dev": {
		Meta: confcenter.Meta{
			Namespace: "dev",
		},
		Server: serverConf,
		Database: confcenter.Database{
			Mongo: mgz.Config{
				Uri:      fmt.Sprintf("mongodb://14.103.49.90:3717/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.MongoPasswordYuzhi,

				Database: "admin",
			},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"118.196.63.209:6379"},
				Password: confs.RedisPasswordProj,
			},

			//Mongo: mgz.Config{
			//	Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
			//	Username: "root",
			//	Password: confs.MongoPasswordProj,
			//	Database: "yoozy_pro_beta",
			//},
			//Rediz: rediz.Config{
			//	Addrs:    []string{"118.196.63.209:6379"},
			//	Password: confs.RedisPasswordProj,
			//	DB:       1,
			//},

			Elastics: elastics.Config{
				//Addresses: []string{"http://localhost:9200/"},
				//Addresses: []string{"http://118.196.63.209:9200/"},
				Addresses: []string{"http://yoozy-qv7.public.cn-hangzhou.es-serverless.aliyuncs.com:9200/"},
				Username:  "yoozy-c86",
				Password:  confs.ElasticPassword,
			},

			Oss: alioss.Config{
				AccessKey:    confs.AliyunOssAccessKey,
				AccessSecret: confs.AliyunOssSecret,

				Bucket:   "yoozy",
				Region:   "oss-cn-hangzhou",
				Endpoint: "oss-cn-hangzhou.aliyuncs.com",
			},
			Tos: tos.Config{
				AccessKeyID:     confs.BytedanceAccessKeyID,
				AccessKeySecret: confs.BytedanceSecretAccessKey,

				//Endpoint:        "tos-s3-cn-shanghai.volces.com",
				Endpoint: "tos-cn-shanghai.volces.com",
				//Endpoint: "tos-accelerate.volces.com",
				//yoozyres.tos-accelerate.volces.com
				DefaultBucket: "yoozyres",
				Region:        "cn-shanghai",
			},
		},
		Component: confcenter.Component{
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
			OpenAI: openaiz.Config{
				AppKey: confs.OpenAIKeys[0],

				//Proxy:   "http://proxy:strOngPAssWOrd@45.78.194.147:6060",
				BaseUrl: "http://45.78.194.147:6000",
			},
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,

				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				ProjAdmin: &grpcz.Config{
					Endpoint: "admin.prod:8090",
				},
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter.prod:8090",
				},
				Credit: &grpcz.Config{
					Endpoint: "usercenter.prod:8090",
				},
			},
		},
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

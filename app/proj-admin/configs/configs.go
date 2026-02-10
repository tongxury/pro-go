package configs

import (
	"fmt"
	"os"
	"store/confs"
	"store/pkg/clients/elastics"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/aliyun/alioss"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

var serverConf = confcenter.Server{
	Http: &confcenter.ServerConfig{
		Addr:    "0.0.0.0:8082",
		Timeout: 30000000000,
	},
	Grpc: &confcenter.ServerConfig{
		Addr:    "0.0.0.0:8092",
		Timeout: 30000000000,
	},
}

type BizConfig struct {
}

var logger = log.With(log.NewStdLogger(os.Stdout),
	"ts", log.DefaultTimestamp,
	"caller", log.DefaultCaller,
	"trace_id", tracing.TraceID(),
	"span_id", tracing.SpanID(),
)

var configs = map[string]*confcenter.Config[BizConfig]{
	"prod": {
		//Logger: logger,
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
			//Clickhouse: confcenter.ClickHouseConfig{
			//	//Addrs:    []string{"95.217.42.20:9000"},
			//	//Database: "default",
			//	//Username: "default",
			//	//Password: "erGVO4a0QL",
			//
			//	//Addrs:    []string{"95.217.42.20:32589"},
			//	//Database: "default",
			//	//Username: "default",
			//
			//	Addrs:    []string{"173.208.218.161:9000"},
			//	Database: "default",
			//	Username: "default",
			//},
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"118.196.63.209:6379"},
				Password: confs.RedisPasswordProj,
			},
			Mongo: mgz.Config{
				//Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Password: confs.MongoPasswordProj,
				//Database: "yoozy_pro",
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
			Oss: alioss.Config{
				AccessKey:    confs.AliyunOssAccessKey,
				AccessSecret: confs.AliyunOssSecret,
				Bucket:       "yoozy",
				Region:       "oss-cn-hangzhou",
				Endpoint:     "oss-cn-hangzhou.aliyuncs.com",
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

					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},

					//{
					//	BaseURL: "https://api.viviai.cc",
					//},
				},
			},
			Grpc: grpcz.Configs{
				ProjPro: &grpcz.Config{
					Endpoint: "proj-pro:8090",
				},
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
				Credit: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
			},
		},
	},
	"beta": {
		//Logger: logger,
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
			//Clickhouse: confcenter.ClickHouseConfig{
			//	//Addrs:    []string{"95.217.42.20:9000"},
			//	//Database: "default",
			//	//Username: "default",
			//	//Password: "erGVO4a0QL",
			//
			//	//Addrs:    []string{"95.217.42.20:32589"},
			//	//Database: "default",
			//	//Username: "default",
			//
			//	Addrs:    []string{"173.208.218.161:9000"},
			//	Database: "default",
			//	Username: "default",
			//},
			Mongo: mgz.Config{
				//Uri:           "173.208.218.161:27017",
				Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.MongoPasswordProj,
				Database: "yoozy_pro_beta",

				//Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "admin",
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
				Bucket:       "yoozy",
				Region:       "oss-cn-hangzhou",
				Endpoint:     "oss-cn-hangzhou.aliyuncs.com",
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

					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},
				},
			},
			Grpc: grpcz.Configs{
				ProjPro: &grpcz.Config{
					Endpoint: "proj-pro:8090",
				},
				UserCenter: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
				Credit: &grpcz.Config{
					Endpoint: "usercenter:8090",
				},
			},
		},
	},
	"dev": {
		//Logger: logger,
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
			//Rediz: rediz.Config{
			//	//Addrs:        []string{"redis-master.prod:6379"},
			//	//Password:     "lveRN3bj7b",
			//	Addrs:    []string{"118.196.63.209:6379"},
			//	Password: confs.RedisPasswordProj,
			//},
			//Mongo: mgz.Config{
			//	//Uri:           "173.208.218.161:27017",
			//	Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
			//	Username: "root",
			//	Password: confs.MongoPasswordProj,
			//	Database: "yoozy_pro_beta",
			//},
			Elastics: elastics.Config{
				//Addresses: []string{"http://localhost:9200/"},
				//Addresses: []string{"http://118.196.63.209:9200/"},
				Addresses: []string{"http://yoozy-qv7.public.cn-hangzhou.es-serverless.aliyuncs.com:9200/"},
				Username:  "yoozy-qv7",
				Password:  confs.ElasticPassword,
			},
			Oss: alioss.Config{
				AccessKey:    confs.AliyunOssAccessKey,
				AccessSecret: confs.AliyunOssSecret,
				Bucket:       "yoozy",
				Region:       "oss-cn-hangzhou",
				Endpoint:     "oss-cn-hangzhou.aliyuncs.com",
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

			//VITE_OSS_BUCKET = yoozy
			//VITE_OSS_REGION = oss-cn-hangzhou
		},
		Component: confcenter.Component{
			Genai: gemini.FactoryConfig{
				Configs: []*gemini.Config{

					{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},

					//{
					//	BaseURL: "https://api.viviai.cc",
					//},
				},
			},
			Grpc: grpcz.Configs{
				ProjPro: &grpcz.Config{
					Endpoint: "proj-pro.prod:8090",
				},
				UserCenter: &grpcz.Config{
					Endpoint: "localhost:8091",
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
		panic("no config found")
	}

	fmt.Println(cc)

	return cc, nil
}

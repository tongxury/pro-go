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
	"store/pkg/sdk/third/aliyun/alisms"
	"store/pkg/sdk/third/bytedance/tos"
)

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
				//Uri:           "173.208.218.161:27017",
				//Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "yoozy_pro",

				Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.MongoPasswordYuzhi,
				Database: "admin",

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
		},
		Component: confcenter.Component{},
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
				//Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				//Username: "root",
				//Database: "yoozy_pro_beta",

				Uri:      fmt.Sprintf("mongodb://mongoreplicad78433cc0e940.mongodb.cn-shanghai.ivolces.com:3717/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.MongoPasswordYuzhi,
				Database: "admin",
			},
			Elastics: elastics.Config{
				//Addresses: []string{"http://localhost:9200/"},
				//Addresses: []string{"http://118.196.63.209:9200/"},
				Addresses: []string{"http://yoozy-qv7.public.cn-hangzhou.es-serverless.aliyuncs.com:9200/"},
				Username:  "yoozy-c86",
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
			},
		},
	},
	"dev": {
		Meta: confcenter.Meta{
			Namespace: "dev",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8071",
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8091",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mongo: mgz.Config{
				//Uri:           "173.208.218.161:27017",
				Uri:      fmt.Sprintf("mongodb://118.196.63.209:27017/yoozy_pro?retryWrites=true&w=majority"),
				Username: "root",
				Password: confs.RedisPasswordProj,

				Database: "yoozy_pro_beta",
			},
			Elastics: elastics.Config{
				//Addresses: []string{"http://localhost:9200/"},
				//Addresses: []string{"http://118.196.63.209:9200/"},
				Addresses: []string{"http://yoozy-qv7.public.cn-hangzhou.es-serverless.aliyuncs.com:9200/"},
				Username:  "yoozy-c86",
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
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,

				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				ProjAdmin: &grpcz.Config{
					Endpoint: "admin.prod.svc.cluster.local:8090",
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

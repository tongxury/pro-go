package configs

import (
	"fmt"
	"os"
	"store/confs"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
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
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"118.196.63.209:6379"},
				Password: confs.RedisPasswordProj,
			},
		},
		Component: confcenter.Component{},
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

			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{"118.196.63.209:6379"},
				Password: confs.RedisPasswordProj,
			},
		},
		Component: confcenter.Component{},
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

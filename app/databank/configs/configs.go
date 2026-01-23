package configs

import (
	"fmt"
	"os"
	"store/app/databank/internal/conf"
	"store/confs"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"

	"entgo.io/ent/dialect"
)

var serverConf = confcenter.Server{
	Http: &confcenter.ServerConfig{
		Addr:    "0.0.0.0:8080",
		Timeout: 30000000000,
	},
	Grpc: &confcenter.ServerConfig{
		Addr:    "0.0.0.0:8090",
		Timeout: 30000000000,
	},
}

var ip = "173.208.218.161"

var configs = map[string]*confcenter.Config[conf.BizConfig]{
	"prod": {
		Meta: confcenter.Meta{
			Namespace: "prod",
		},
		Server: serverConf,
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordProj, ip),
				//Driver: dialect.Postgres,
			},
			Clickhouse: confcenter.ClickHouseConfig{
				//Addrs:    []string{"clickhouse.prod:9000"},
				//Database: "default",

				//Username: "default",
				//Password: "erGVO4a0QL",

				//Addrs:    []string{"clickhouse.prod:9000"},
				//Database: "default",
				//Username: "default",
				//Password: confs.ClickhousePassword,

				Addrs:    []string{ip + ":9000"},
				Database: "default",
				Username: "default",
				Password: confs.ClickhousePassword,
			},
			Mongo: mgz.Config{
				//Uri:      "mongodb://mongodb-sharded.prod:27017/tgbot?retryWrites=true&w=majority",
				//Username: "root",
				//Password: "z4XNmlaOjo",
				Uri:      fmt.Sprintf("mongodb://%s:27017/pro?retryWrites=true&w=majority", ip),
				Username: "root",
				Password: confs.MongoPassword,
				Database: "pro",
			},
			Redis: confcenter.Redis{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:        []string{ip + ":16379"},
				Password:     confs.RedisPassword,
				ReadTimeout:  10000000000,
				WriteTimeout: 10000000000,
			},
		},
		Component: confcenter.Component{
			Kafka: confcenter.KafkaConfig{
				Brokers: []string{"kafka-headless.prod.svc.cluster.local:9092"},
				//Brokers: []string{
				//	"kafka-controller-0.kafka-controller-headless.prod.svc.cluster.local:9092",
				//	"kafka-controller-1.kafka-controller-headless.prod.svc.cluster.local:9092",
				//	"kafka-controller-2.kafka-controller-headless.prod.svc.cluster.local:9092",
				//},
				//Username: "user1",
				//Password: "QpF5Y1x0M7",
			},
			Solana: confcenter.SolanaConfig{
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			Grpc: grpcz.Configs{
				TgBot: &grpcz.Config{
					Endpoint: "dex-tgbot-server.prod.svc.cluster.local:8090",
				},
				TgUser: &grpcz.Config{
					Endpoint: "dex-tgbot-server.prod.svc.cluster.local:8090",
				},
				DexTrade: &grpcz.Config{
					Endpoint: "dex-tgbot-server.prod.svc.cluster.local:8090",
				},
				DexMarket: &grpcz.Config{
					Endpoint: "dex-tgbot-server.prod.svc.cluster.local:8090",
				},
				DexWallet: &grpcz.Config{
					Endpoint: "dex-tgbot-server.prod.svc.cluster.local:8090",
				},
				DexAdmin: &grpcz.Config{
					Endpoint: "dex-admin.prod.svc.cluster.local:8090",
				},
			},
		},
		Biz: conf.BizConfig{},
	},
	"dev": {
		Meta: confcenter.Meta{
			Namespace: "dev",
		},
		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8083",
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8093",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordProj, ip),
				//Driver: dialect.Postgres,
			},
			Clickhouse: confcenter.ClickHouseConfig{
				//Addrs:    []string{"95.217.42.20:9000"},
				//Database: "default",
				//Username: "default",
				//Password: "erGVO4a0QL",

				//Addrs:    []string{"95.217.42.20:32589"},
				//Database: "default",
				//Username: "default",
				//Password: confs.ClickhousePassword,

				Addrs:    []string{ip + ":9000"},
				Database: "default",
				Username: "default",
				Password: confs.ClickhousePassword,
			},
			Mongo: mgz.Config{
				Uri:      "mongodb://" + ip + ":27017/dex?retryWrites=true&w=majority",
				Username: "root",
				Password: confs.MongoPassword,

				//Uri:      "mongodb://localhost:27017/tgbot?retryWrites=true&w=majority",
				//Username: "root",
				//Password: "pRomonGo",

				//Uri:      "mongodb://95.217.42.20:27017/dex?retryWrites=true&w=majority",
				//Username: "root",
				//Password: "z4XNmlaOjo",
				Database: "pro",
			},
			Redis: confcenter.Redis{
				//Addrs: []string{"127.0.0.1:6379"},
				//Password:     confs.RedisPassword,
				Addrs:        []string{ip + ":16379"},
				Password:     confs.RedisPassword,
				ReadTimeout:  10000000000,
				WriteTimeout: 10000000000,
			},
		},
		Component: confcenter.Component{
			Solana: confcenter.SolanaConfig{
				//Endpoint:   "https://dark-radial-seed.solana-mainnet.quiknode.pro/95494d59ce7464c3b374dcae1d25c0a3cba837f2",
				//WSEndpoint: "wss://dark-radial-seed.solana-mainnet.quiknode.pro/95494d59ce7464c3b374dcae1d25c0a3cba837f2",
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			Grpc: grpcz.Configs{
				TgBot: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
				TgUser: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
				DexTrade: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
				DexMarket: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
				DexWallet: &grpcz.Config{
					Endpoint: "localhost:8090",
				},
				DexAdmin: &grpcz.Config{
					Endpoint: "localhost:8091",
				},
			},
		},
		Biz: conf.BizConfig{},
	},
}

func GetConfig() (*confcenter.Config[conf.BizConfig], error) {

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

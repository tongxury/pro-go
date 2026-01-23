package configs

import (
	"fmt"
	"os"
	bizconf "store/app/payment/internal/conf"
	"store/confs"
	"store/pkg/clients/grpcz"

	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/airwallex"
	"store/pkg/sdk/third/aws"
	"store/pkg/sdk/third/aws/awsmail"
	"store/pkg/sdk/third/stripe"

	"entgo.io/ent/dialect"
	"github.com/go-kratos/kratos/v2/log"
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

var ip = "13.250.123.59"
var hkip = "47.76.179.31"

var configs = map[string]*confcenter.Config[bizconf.BizConfig]{
	"prod-cn": {
		Meta: confcenter.Meta{
			Namespace: "prod",
		},
		Server: serverConf,
		Logger: log.DefaultLogger,
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordProj, "101.132.192.41"),
				//Driver: dialect.Postgres,
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
		},
		Component: confcenter.Component{
			Kafka: confcenter.KafkaConfig{
				Brokers: []string{"101.132.192.41:9092"},
				//Brokers: []string{
				//	"kafka-controller-0.kafka-controller-headless.prod.svc.cluster.local:9092",
				//	"kafka-controller-1.kafka-controller-headless.prod.svc.cluster.local:9092",
				//	"kafka-controller-2.kafka-controller-headless.prod.svc.cluster.local:9092",
				//},
				//Username: "user1",
				//Password: "QpF5Y1x0M7",
			},
			AirWallex: airwallex.Config{
				ClientId:     confs.AirWallexClientId,
				ClientSecret: confs.AirWallexClientSecret,

				//ClientId:     "9jtRbZaXSAy5GcPGgZVHeA",
				//ClientSecret: "9fd7512838865b660d0c4de708b8edb28282b93a633a335bfd05eeca22567e7e17eea9dc22487cb2d687be77d0184178",
				Endpoint:       "https://api.airwallex.com",
				AccountId:      "acct_Q6NUqVssMGu--2O8trlsVw",
				CallbackSecret: confs.AirWallexCallbackSecret,

			},
			StripePay: confcenter.StripePay{
				Prod: stripe.Config{
					Key:    confs.StripeKeyLive,
					Secret: confs.StripeSecretLive,

					BillLoginUrl: "",
					Emails:       nil,
				},
				Test: stripe.Config{
					Key:    confs.StripeKeyTestLegacy,
					Secret: confs.StripeKeyTest,
				},
			},

			Grpc: grpcz.Configs{
				User: &grpcz.Config{
					Endpoint: "user.prod.svc.cluster.local:8090",
				},
			},
		},
		//Biz: conf.BizConfig{},
	},
	"prod": {
		Meta: confcenter.Meta{
			Namespace: "prod",
			Domain:    "https://i.veogo.ai",
		},
		Server: serverConf,
		Logger: log.DefaultLogger,
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
			Rediz: rediz.Config{
				//Addrs:        []string{"redis-master.prod:6379"},
				//Password:     "lveRN3bj7b",
				Addrs:    []string{ip + ":16379"},
				Password: confs.RedisPassword,

				//ReadTimeout:  10000000000,
				//WriteTimeout: 10000000000,
			},
		},
		Component: confcenter.Component{
			AirWallex: airwallex.Config{
				ClientId:     confs.AirWallexClientId,
				ClientSecret: confs.AirWallexClientSecret,

				//ClientId:     "9jtRbZaXSAy5GcPGgZVHeA",
				//ClientSecret: "9fd7512838865b660d0c4de708b8edb28282b93a633a335bfd05eeca22567e7e17eea9dc22487cb2d687be77d0184178",
				Endpoint:  "https://api.airwallex.com",
				AccountId: "acct_Q6NUqVssMGu--2O8trlsVw",
				CallbackSecret: confs.AirWallexCallbackSecretLegacy,
			},
			StripePay: confcenter.StripePay{
				Prod: stripe.Config{
					Key:    confs.StripeKeyLive,
					Secret: confs.StripeSecretLive,

					BillLoginUrl: "",
					Emails:       nil,
				},
				Test: stripe.Config{
					Key:    confs.StripeKeyTestLegacy,
					Secret: confs.StripeKeyTest,
				},
			},
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
			AwsMail: awsmail.Config{
				AwsConfig: aws.AwsConfig{
					AccessKey:    confs.AWSAccessKey,
					AccessSecret: confs.AWSAccessSecret,

					Region: "us-east-2",
				},
				Sender: "noreply@veogo.ai",
			},
			Grpc: grpcz.Configs{
				User: &grpcz.Config{
					Endpoint: "user.prod.svc.cluster.local:8090",
				},
			},
		},
	},
	"dev": {
		Meta: confcenter.Meta{
			Namespace: "dev",
			Domain:    "https://i.veogo.ai",
		},
		Logger: log.DefaultLogger,

		Server: confcenter.Server{
			Http: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8085",
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8095",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordProj, "13.250.123.59"),
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
				Addrs:    []string{ip + ":16379"},
				Password: confs.RedisPassword,

				ReadTimeout:  10000000000,
				WriteTimeout: 10000000000,
			},
		},
		Component: confcenter.Component{
			StripePay: confcenter.StripePay{
				Prod: stripe.Config{
					Key:    confs.StripeKeyLive,
					Secret: confs.StripeSecretLive,

					BillLoginUrl: "",
					Emails:       nil,
				},
			},
			Solana: confcenter.SolanaConfig{
				//Endpoint:   "https://dark-radial-seed.solana-mainnet.quiknode.pro/95494d59ce7464c3b374dcae1d25c0a3cba837f2",
				//WSEndpoint: "wss://dark-radial-seed.solana-mainnet.quiknode.pro/95494d59ce7464c3b374dcae1d25c0a3cba837f2",
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			AwsMail: awsmail.Config{
				AwsConfig: aws.AwsConfig{
					AccessKey:    confs.AWSAccessKey,
					AccessSecret: confs.AWSAccessSecret,

					Region: "us-east-2",
				},
				Sender: "noreply@veogo.ai",
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
		Biz: bizconf.BizConfig{},
	},
}

func GetConfig() (*confcenter.Config[bizconf.BizConfig], error) {

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

package configs

import (
	"fmt"
	"os"
	"store/app/user/internal/conf"
	"store/confs"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/aliyun/alisms"
	"store/pkg/sdk/third/aws"
	"store/pkg/sdk/third/aws/awsmail"

	"entgo.io/ent/dialect"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/oauth2"
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

var ip = "101.132.192.41"
var hkip = "47.76.179.31"

var configs = map[string]*confcenter.Config[conf.BizConfig]{
	"proj-prod": {
		Meta: confcenter.Meta{
			Namespace: "prod",
			Domain:    "https://i.veogo.cn",
		},
		Logger: log.DefaultLogger,
		Server: serverConf,
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(101.132.192.41:3306)/yoozy?parseTime=True&loc=Local", confs.MysqlPasswordProj),
				//Driver: dialect.Postgres,
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
				Brokers: []string{"54.222.192.154:9092"},
			},
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,
				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				AiAgent: &grpcz.Config{
					Endpoint: "aiagent.prod.svc.cluster.local:8090",
				},
				Payment: &grpcz.Config{
					Endpoint: "payment.prod.svc.cluster.local:8090",
				},
			},
		},
		Biz: conf.BizConfig{
			Oauth2: conf.Oauth2{
				Redirect:     nil,
				CookieDomain: "veogo.ai",
				Google: oauth2.Config{
					ClientID:     confs.GoogleClientID,
					ClientSecret: confs.GoogleClientSecret,
				},
			},
		},
	},
	"prod-cn": {
		Meta: confcenter.Meta{
			Namespace: "prod",
			Domain:    "https://i.veogo.cn",
		},
		Logger: log.DefaultLogger,
		Server: serverConf,
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("admin:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordAdmin, "db.cfqaeqikijtq.rds.cn-north-1.amazonaws.com.cn"),
			},
			Clickhouse: confcenter.ClickHouseConfig{
				Addrs:    []string{ip + ":9000"},
				Database: "default",
				Username: "default",
				Password: confs.ClickhousePassword,
			},
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
		},
		Component: confcenter.Component{
			Kafka: confcenter.KafkaConfig{
				Brokers: []string{"101.132.192.41:9092"},
			},
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,
				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				AiAgent: &grpcz.Config{
					Endpoint: "aiagent.prod.svc.cluster.local:8090",
				},
				Payment: &grpcz.Config{
					Endpoint: "payment.prod.svc.cluster.local:8090",
				},
			},
		},
		Biz: conf.BizConfig{
			Oauth2: conf.Oauth2{
				Redirect:     nil,
				CookieDomain: "veogo.ai",
				Google: oauth2.Config{
					ClientID:     confs.GoogleClientID,
					ClientSecret: confs.GoogleClientSecret,
				},
			},
		},
	},
	"prod": {
		Meta: confcenter.Meta{
			Namespace: "prod",
			Domain:    "https://i.veogo.ai",
		},
		Logger: log.DefaultLogger,
		Server: serverConf,
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
				Uri:      fmt.Sprintf("mongodb://%s:27017/pro?retryWrites=true&w=majority", ip),
				Username: "root",
				Password: confs.MongoPassword,
				Database: "pro",
			},
			Rediz: rediz.Config{
				Addrs:    []string{ip + ":6379"},
				Password: confs.RedisPassword,
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
			AwsMail: awsmail.Config{
				AwsConfig: aws.AwsConfig{
					AccessKey:    confs.AWSAccessKey,
					AccessSecret: confs.AWSAccessSecret,
					Region:       "us-east-2",
				},
				Sender: "noreply@veogo.ai",
			},
			Alisms: alisms.Config{
				AccessKey:    confs.AliyunAccessKey,
				AccessSecret: confs.AliyunAccessSecret,
				Sign:         "唯构科技深圳",
				TemplateCode: "SMS_316000047",
			},
			Grpc: grpcz.Configs{
				AiAgent: &grpcz.Config{
					Endpoint: "aiagent.prod.svc.cluster.local:8090",
				},
				Payment: &grpcz.Config{
					Endpoint: "payment.prod.svc.cluster.local:8090",
				},
			},
		},
		Biz: conf.BizConfig{
			Oauth2: conf.Oauth2{
				Redirect:     nil,
				CookieDomain: "veogo.ai",
				Google: oauth2.Config{
					ClientID:     confs.GoogleClientID,
					ClientSecret: confs.GoogleClientSecret,
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
				Addr:    "0.0.0.0:8088",
				Timeout: 30000000000,
			},
			Grpc: &confcenter.ServerConfig{
				Addr:    "0.0.0.0:8098",
				Timeout: 30000000000,
			},
		},
		Database: confcenter.Database{
			Mysql: confcenter.Mysql{
				Driver: dialect.MySQL,
				Source: fmt.Sprintf("root:%s@tcp(%s:3306)/pro?parseTime=True&loc=Local", confs.MysqlPasswordProj, "13.250.123.59"),
			},
			Clickhouse: confcenter.ClickHouseConfig{
				Addrs:    []string{ip + ":9000"},
				Database: "default",
				Username: "default",
				Password: confs.ClickhousePassword,
			},
			Mongo: mgz.Config{
				Uri:      "mongodb://" + ip + ":27017/dex?retryWrites=true&w=majority",
				Username: "root",
				Password: confs.MongoPassword,
				Database: "pro",
			},
			Rediz: rediz.Config{
				Addrs:    []string{ip + ":6379"},
				Password: confs.RedisPassword,
			},
		},
		Component: confcenter.Component{
			Solana: confcenter.SolanaConfig{
				Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
				WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
			},
			AwsMail: awsmail.Config{
				AwsConfig: aws.AwsConfig{
					AccessKey:    confs.AWSAccessKey,
					AccessSecret: confs.AWSAccessSecret,
					Region:       "us-east-2",
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
		return nil, fmt.Errorf("env %s not found", env)
		//return configs["dev"], nil
	}

	return cc, nil
}

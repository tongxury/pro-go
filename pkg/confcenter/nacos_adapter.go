package confcenter

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
)

type NacosClient struct {
	config_client.IConfigClient
}

func NewClient(env string) (*NacosClient, error) {

	// 配置中心基本不变先写死
	host := "confcenter.components.svc.cluster.local"
	port := 8080
	grpcPort := 8090

	// 默认线上
	if env == "" {
		env = "prod"
	}
	//namespace := "prod"
	//switch env {
	//case "staging":
	//	namespace = "21be335a-281b-449a-8d91-845c6e49e6f4"
	//}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				NamespaceId:         env,
				NotLoadCacheAtStart: true,
				//Username:            "nacos",
				//Password:            "nacos",
				TimeoutMs: 10000,
				LogLevel:  "info",
			},
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig(
					host,
					uint64(port),
					constant.WithGrpcPort(uint64(grpcPort)),
					constant.WithScheme("http"),
					constant.WithContextPath("/nacos"),
				),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	return &NacosClient{
		client,
	}, nil
}

func (t *NacosClient) LoadDatabase(options ...Options) (*Database, error) {

	dataId := "database.yaml"
	group := "public"

	databaseSource, err := t.GetConfig(vo.ConfigParam{
		DataId: dataId, Group: group,
	})

	if err != nil {
		return nil, err
	}

	var database *Database
	err = yaml.Unmarshal([]byte(databaseSource), &database)
	if err != nil {
		return nil, err
	}

	if len(options) > 0 {
		if options[0].Listen {
			err = t.ListenConfig(vo.ConfigParam{
				DataId: dataId, Group: group, OnChange: func(namespace, group, dataId, data string) {
					err = yaml.Unmarshal([]byte(data), &database)
					if options[0].OnChange != nil {
						options[0].OnChange(data)
					}
				},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return database, nil
}

func (t *NacosClient) WatchDatabase() (*Database, error) {
	databaseSource, err := t.GetConfig(vo.ConfigParam{
		DataId: "database.yaml", Group: "public",
	})

	if err != nil {
		return nil, err
	}

	var database Database
	err = yaml.Unmarshal([]byte(databaseSource), &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

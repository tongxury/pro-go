package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"store/app/databank/internal/data/repo/ent"
	"store/pkg/clients"
	"store/pkg/confcenter"
)

func NewEntClient(conf confcenter.Mysql) *ent.Client {

	client, err := ent.Open(conf.Driver, conf.Source)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}

	// todo 容易把字段清掉 还没看原因
	//Run the auto migration tool.
	err = client.Schema.Create(
		context.Background(),
		//migrate.WithDropIndex(true),
		//migrate.WithDropColumn(true),
		//migrate.WithForeignKeys(false),
	)

	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

type Repos struct {
	EntClient   *ent.Client
	RedisClient *clients.RedisClient
	UserAsset   *UserAssetRepo
	UserFile    *UserFileRepo
}

func NewRepos(entClient *ent.Client, redisClient *clients.RedisClient,

) *Repos {

	return &Repos{
		EntClient:   entClient,
		RedisClient: redisClient,
		UserAsset:   NewUserAssetRepo(entClient, redisClient),
		UserFile:    NewUserFileRepo(entClient, redisClient),
	}
}

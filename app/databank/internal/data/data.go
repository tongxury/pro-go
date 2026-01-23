package data

import (
	"store/app/databank/internal/conf"
	"store/app/databank/internal/data/repo"
	"store/pkg/clients"
	"store/pkg/confcenter"
	"store/pkg/sdk/third/aliyun/alioss"
)

// Data .
type Data struct {
	BizConfig conf.BizConfig
	Repos     *repo.Repos
	Uploader  *Uploader
}

func NewData(
	conf confcenter.Config[conf.BizConfig],
) (*Data, func(), error) {

	entClient := repo.NewEntClient(conf.Database.Mysql)
	redisClient := clients.NewRedisClient(conf.Database.Redis)

	d := &Data{
		BizConfig: conf.Biz,
		Repos:     repo.NewRepos(entClient, redisClient),
		Uploader:  NewUploader(alioss.NewClient(conf.Database.Oss)),
	}
	return d, func() {
		//if err := d.DB.Close(); err != nil {
		//	log.Error(err)
		//}
		//if err := d.RedisCluster.Close(); err != nil {
		//	log.Error(err)
		//}
	}, nil
}

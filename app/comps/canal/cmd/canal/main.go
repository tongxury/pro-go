package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"store/app/comps/canal/internal/conf"
	"store/pkg/clients"
	"store/pkg/confcenter"
	"store/pkg/events"
	"store/pkg/sdk/conv"
)

func main() {

	config, err := confcenter.GetConfig[conf.BizConfig]()
	if err != nil {
		log.Fatalf("failed get config: %v", err)
	}

	kafkaClient := clients.NewKafkaClient(config.Component.Kafka)

	dsn, err := mysql.ParseDSN(config.Database.Mysql.Source)
	if err != nil {
		log.Fatalf("failed to ParseDSN: %v", err)
	}

	log.Debugw("mysql config", dsn)

	c := NewCanalClient(CanalConfig{
		Addr:     dsn.Addr,
		User:     dsn.User,
		Password: dsn.Passwd,
		//Addr:      "prod-db.cdetyl3o49s3.us-east-1.rds.amazonaws.com:3306",
		//User:      "admin",
		//Password:  "xbuddy.ai",

		Databases: nil,
		Tables:    []string{},
	})

	c.SetEventHandler(NewEventHandler(func(e *canal.RowsEvent) error {

		ev := events.FromRowsEvent(e)
		log.Debug(e.Rows, ev)

		ctx := context.Background()

		err := kafkaClient.W().Write(ctx, events.Topic_MysqlRowUpdate, kafka.Message{
			Value: conv.S2B(ev),
		})
		if err != nil {
			log.Errorw("Write Kafka err", err, "e", e)
			return err
		}

		return nil
	}))

	err = c.Run()
	if err != nil {
		log.Fatalw("start err", err)
	}
}

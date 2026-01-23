package mgz

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Uri           string
	AuthMechanism string
	Username      string
	Password      string
	Database      string
}

func Database(config Config) (*mongo.Database, error) {
	mongoClient, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(config.Uri).
			SetAuth(options.Credential{
				Username: config.Username,
				Password: config.Password,
			}),
	)

	if err != nil {
		return nil, err
	}

	database := mongoClient.Database(config.Database)

	return database, nil
}

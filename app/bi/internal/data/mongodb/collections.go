package mongodb

import (
	"store/pkg/clients/mgz"
)

type Collections struct {
	JuGuang *JuGuangCollection
}

func NewCollections(config mgz.Config) *Collections {

	database, err := mgz.Database(config)

	if err != nil {
		panic(err)
	}

	return &Collections{
		JuGuang: NewJuGuangCollection(database),
	}
}

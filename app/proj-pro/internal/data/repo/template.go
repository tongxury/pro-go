package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateCollection struct {
	*mgz.Core[*projpb.Resource]
}

func NewTemplateCollection(db *mongo.Database) *TemplateCollection {
	return &TemplateCollection{
		Core: mgz.NewCore[*projpb.Resource](db, "templates"),
	}
}

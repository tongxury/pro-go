package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateSegmentCollection struct {
	*mgz.Core[*projpb.ResourceSegment]
}

func NewTemplateSegmentCollection(db *mongo.Database) *TemplateSegmentCollection {
	return &TemplateSegmentCollection{
		Core: mgz.NewCore[*projpb.ResourceSegment](db, "template_segments"),
	}
}

package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceSegmentCollection struct {
	*mgz.Core[*projpb.ResourceSegmentCollectionItem]
}

func NewResourceSegmentCollection(db *mongo.Database) *ResourceSegmentCollection {
	return &ResourceSegmentCollection{
		Core: mgz.NewCore[*projpb.ResourceSegmentCollectionItem](db, "template_segment_collection"),
	}
}

package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
)

type TaskSegmentCollection struct {
	*mgz.CoreCollectionV3[projpb.TaskSegment]
}

func NewTaskSegmentCollection(db *mongo.Database) *TaskSegmentCollection {
	return &TaskSegmentCollection{
		CoreCollectionV3: mgz.NewCoreCollectionV3[projpb.TaskSegment](db, "task_segments"),
	}
}

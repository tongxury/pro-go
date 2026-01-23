package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type TaskCollection struct {
	*mgz.Core[*projpb.Task]
}

func NewTaskCollection(db *mongo.Database) *TaskCollection {
	return &TaskCollection{
		Core: mgz.NewCore[*projpb.Task](db, "tasks"),
	}
}

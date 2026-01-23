package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowCollection struct {
	*mgz.Core[*projpb.Workflow]
}

func NewWorkflowCollection(db *mongo.Database) *WorkflowCollection {
	return &WorkflowCollection{
		Core: mgz.NewCore[*projpb.Workflow](db, "workflows"),
	}
}

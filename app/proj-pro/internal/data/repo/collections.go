package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Task            *TaskCollection
	TaskSegment     *TaskSegmentCollection
	Commodity       *CommodityCollection
	Settings        *SettingsCollection
	Asset           *AssetCollection
	Template        *TemplateCollection
	TemplateSegment *TemplateSegmentCollection
	Feedback        *FeedbackCollection
	Session         *SessionCollection
	SessionSegment  *SessionSegmentCollection
	Workflow        *WorkflowCollection
}

func NewCollections(database *mongo.Database) *Collections {
	return &Collections{
		Task:            NewTaskCollection(database),
		Commodity:       NewCommodityCollection(database),
		Settings:        NewSettingsCollection(database),
		TaskSegment:     NewTaskSegmentCollection(database),
		Asset:           NewAssetCollection(database),
		Template:        NewTemplateCollection(database),
		TemplateSegment: NewTemplateSegmentCollection(database),
		Feedback:        NewFeedbackCollection(database),
		Session:         NewSessionCollection(database),
		SessionSegment:  NewSessionSegmentCollection(database),
		Workflow:        NewWorkflowCollection(database),
	}
}

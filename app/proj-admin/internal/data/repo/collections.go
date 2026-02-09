package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Template        *TemplateCollection
	TemplateSegment *TemplateSegmentCollection
	Settings        *SettingsCollection
	Database        *mongo.Database
}

func NewCollections(database *mongo.Database) *Collections {
	return &Collections{
		Template:        NewTemplateCollection(database),
		TemplateSegment: NewTemplateSegmentCollection(database),
		Settings:        NewSettingsCollection(database),
		Database:        database,
	}
}

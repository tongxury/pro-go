package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"store/pkg/enums"
	"store/pkg/middlewares/entz"
	"store/pkg/types"
)

type Notification struct {
	ent.Schema
}

func (Notification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "notifications"},
	}
}

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		entz.CreatedAtField(),
		field.Enum("level").GoType(enums.NotificationLevel_Info),
		entz.TimeField("start_at"),
		entz.TimeField("end_at"),
		entz.ExtraField(types.NotificationExtra{}),
	}
}

func (Notification) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("start_at"),
		index.Fields("end_at"),
	}
}

func (Notification) Edges() []ent.Edge {
	return nil
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"store/pkg/middlewares/entz"
)

// User holds the schema definition for the User entity.
type Feedback struct {
	ent.Schema
}

func (Feedback) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "feedbacks"},
	}
}

// Fields of the User.
func (Feedback) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		field.Int64("user_id").Immutable(),
		field.String("category").NotEmpty().Immutable(),
		field.Text("content").NotEmpty().Immutable(),
		entz.CreatedAtField(),
		entz.ExtraField(map[string]any{}),
	}
}

func (Feedback) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}

// Edges of the User.
func (Feedback) Edges() []ent.Edge {
	return nil
}

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
type UserFile struct {
	ent.Schema
}

func (UserFile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_files"},
	}
}

// Fields of the User.
func (UserFile) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		field.Int64("user_id").Immutable(),
		field.String("name").NotEmpty().Immutable(),
		field.String("md5").NotEmpty().Immutable(),
		field.Int64("size").Immutable(),
		field.String("category"),
		entz.CreatedAtField(),
		entz.ExtraField(map[string]any{}),
	}
}

func (UserFile) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "name"),
		index.Fields("user_id", "md5").Unique(),
	}
}

// Edges of the User.
func (UserFile) Edges() []ent.Edge {
	return nil
}

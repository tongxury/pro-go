package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"store/app/user/internal/data/enums"
	"store/pkg/middlewares/entz"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users"},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		entz.CreatedAtField(),
		field.String("key").Unique(),
		field.String("nickname").Optional(),
		field.String("email").Optional(),
		field.String("phone").Optional(),
		field.String("password").Optional(),
		field.String("status").Default(enums.UserStatusPending),
		field.String("channel").Optional(),
		field.String("avatar").Optional(),
		field.String("platform").Optional(),
		//field.String("inviterUserId").Optional(),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("phone"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

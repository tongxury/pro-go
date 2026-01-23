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
type UserCreditIncrement struct {
	ent.Schema
}

func (UserCreditIncrement) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_credit_increments"},
	}
}

// Fields of the User.
func (UserCreditIncrement) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		entz.CreatedAtField(),
		field.String("user_id"),
		field.Int64("amount"),
		entz.TimeField("expire_at"),
		field.String("status").Default(enums.UserCreditStatusPending),
		field.String("plan_id").Optional(),
		field.String("by").Optional(),
	}
}

func (UserCreditIncrement) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}

// Edges of the User.
func (UserCreditIncrement) Edges() []ent.Edge {
	return nil
}

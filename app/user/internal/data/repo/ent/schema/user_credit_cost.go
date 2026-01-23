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
type UserCreditCost struct {
	ent.Schema
}

func (UserCreditCost) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_credit_costs"},
	}
}

// Fields of the User.
func (UserCreditCost) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		entz.CreatedAtField(),
		field.String("user_id"),
		field.Int64("amount"),
	}
}

func (UserCreditCost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}

// Edges of the User.
func (UserCreditCost) Edges() []ent.Edge {
	return nil
}

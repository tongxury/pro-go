package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"store/pkg/middlewares/entz"
)

type CreditConsume struct {
	ent.Schema
}

func (CreditConsume) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "credit_consumes"},
	}
}

// Fields of the User.
func (CreditConsume) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		entz.CreatedAtField(),
		field.String("key").Unique(),
		field.String("recharge_id"),
		field.String("user_id"),
		field.Int64("amount"),
	}
}

func (CreditConsume) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}

// Edges of the User.
func (CreditConsume) Edges() []ent.Edge {
	return nil
}

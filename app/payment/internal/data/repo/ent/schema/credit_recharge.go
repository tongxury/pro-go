package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"store/pkg/middlewares/entz"
)

type CreditRecharge struct {
	ent.Schema
}

func (CreditRecharge) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "credit_recharges"},
	}
}

// Fields of the User.
func (CreditRecharge) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		entz.CreatedAtField(),
		field.String("key").Unique(),
		field.String("user_id"),
		field.String("plan_id"),
		field.Int64("quota"),
		field.Int64("cost").Default(0),
		field.String("status").Default("pending"),
		field.String("platform").Optional(),
		field.Time("expire_at").Optional(),
	}
}

func (CreditRecharge) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "expire_at"),
	}
}

// Edges of the User.
func (CreditRecharge) Edges() []ent.Edge {
	return nil
}

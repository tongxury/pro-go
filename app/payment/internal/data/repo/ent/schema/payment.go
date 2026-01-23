package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"store/pkg/enums"
	"store/pkg/middlewares/entz"
)

// User holds the schema definition for the User entity.
type Payment struct {
	ent.Schema
}

func (Payment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_payments"},
	}
}

// Fields of the User.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		field.String("session_id").Optional().Unique(),
		field.String("sub_id").Optional(),
		field.String("platform").Optional(),
		field.String("user_id").NotEmpty(),
		field.String("plan_id").NotEmpty(),
		field.Float("amount"),
		field.Time("expire_at").Optional(),
		field.Enum("status").GoType(enums.PaymentStatus_Created).Default(enums.PaymentStatus_Created.String()),
		entz.CreatedAtField(),
		entz.ExtraField(map[string]any{}),
	}
}

func (Payment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform", "sub_id"),
		index.Fields("user_id", "expire_at"),
	}
}

// Edges of the User.
func (Payment) Edges() []ent.Edge {
	return nil
}

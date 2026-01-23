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

type Promotion struct {
	ent.Schema
}

func (Promotion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotions"},
	}
}

func (Promotion) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		field.String("code"),
		field.Enum("level").GoType(enums.MemberLevel_Free),
		field.Enum("cycle").GoType(enums.PaymentCycle_Monthly),
		field.Int64("limit"),
		entz.TimeField("start_at"),
		entz.TimeField("end_at"),
		entz.CreatedAtField(),
	}
}

func (Promotion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code"),
	}
}

func (Promotion) Edges() []ent.Edge {
	return nil
}

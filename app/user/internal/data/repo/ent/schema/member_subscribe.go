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

type MemberSubscribe struct {
	ent.Schema
}

func (MemberSubscribe) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_member_subscribes"},
	}
}

func (MemberSubscribe) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id"),
		entz.CreatedAtField(),
		field.Enum("level").GoType(enums.MemberLevel_Free),
		field.Enum("status").GoType(enums.MemberSubscribeStatus_Pending),
		field.Enum("payment_cycle").GoType(enums.PaymentCycle_Monthly),
		field.String("promotion_code").Optional(),
		entz.TimeField("expire_at"),
		entz.ExtraField(types.MemberSubscribeExtra{}),
		field.String("uk"),
	}
}

func (MemberSubscribe) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("uk").Unique(),
	}
}

func (MemberSubscribe) Edges() []ent.Edge {
	return nil
}

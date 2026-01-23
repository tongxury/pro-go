package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"store/pkg/enums"
	"store/pkg/middlewares/entz"
)

// User holds the schema definition for the User entity.
type UserAsset struct {
	ent.Schema
}

func (UserAsset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_assets"},
	}
}

// Fields of the User.
func (UserAsset) Fields() []ent.Field {
	return []ent.Field{
		entz.IdField(),
		field.Int64("user_id").Immutable(),
		field.String("name").NotEmpty().Immutable(),
		//field.String("url").NotEmpty(),
		field.String("key").NotEmpty().Unique().Immutable(),
		field.String("school").Optional(),
		field.String("course").Optional(),
		field.String("year").Optional(),
		field.Enum("status").GoType(enums.AssetStatus_Pending).Default(enums.AssetStatus_Pending.String()),
		field.Enum("category").GoType(enums.AssetCategory_File).Default(enums.AssetCategory_File.String()),
		entz.CreatedAtField(),
		entz.ExtraField(map[string]any{}),
	}
}

func (UserAsset) Indexes() []ent.Index {
	return []ent.Index{
		//index.Fields("phone").Unique(),
	}
}

// Edges of the User.
func (UserAsset) Edges() []ent.Edge {
	return nil
}

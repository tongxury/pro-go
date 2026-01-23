package entz

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"time"
)

func IdField() ent.Field {
	return field.Int64("id")
}

func TimeField(fieldName string, defaultFn ...any) ent.Field {
	f := field.Time(fieldName).Optional()

	if len(defaultFn) > 0 {
		f = f.Default(defaultFn[0])
	}

	f = f.SchemaType(map[string]string{dialect.MySQL: "timestamp", dialect.Postgres: "timestamp"})

	return f
}

func CreatedAtField() ent.Field {
	return field.Time("created_at").Default(time.Now).
		SchemaType(map[string]string{dialect.MySQL: "timestamp", dialect.Postgres: "timestamp"}).
		Immutable()
}

func LastUpdatedAtField() ent.Field {
	return field.Time("last_updated_at").Default(time.Now).
		SchemaType(map[string]string{dialect.MySQL: "timestamp", dialect.Postgres: "timestamp"})
}

func AmountField(fieldName string) ent.Field {
	return field.Float(fieldName).Min(0).
		SchemaType(map[string]string{dialect.MySQL: "decimal(6,2)", dialect.Postgres: "numeric"}).
		Immutable()
}

func ExtraField(typ any) ent.Field {
	return field.JSON("extra", typ).Default(typ)
}

func JsonField(fieldName string, typ any) ent.Field {
	return field.JSON(fieldName, typ).Default(typ).Optional()
}

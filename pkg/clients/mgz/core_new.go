package mgz

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IColl interface {
	GetID() string
	SetID(id string)
}

type Core[T IColl] struct {
	c  *mongo.Collection
	db *mongo.Database
}

func NewCore[T IColl](db *mongo.Database, name string) *Core[T] {
	return &Core[T]{
		db: db,
		c: db.Collection(name, options.Collection().
			SetBSONOptions(&options.BSONOptions{UseJSONStructTags: true})),
	}
}

func (t *Core[T]) GetById(ctx context.Context, id string, options ...*options.FindOptions) (T, error) {
	return t.FindByID(ctx, id, options...)
}

func (t *Core[T]) FindByID(ctx context.Context, id string, options ...*options.FindOptions) (T, error) {
	var zero T
	var ts []T

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return zero, err
	}
	cursor, err := t.c.Find(ctx, bson.M{"_id": hex}, options...)
	if err != nil {
		return zero, err
	}

	err = cursor.All(ctx, &ts)
	if err != nil {
		return zero, err
	}

	if len(ts) == 0 {
		return zero, nil
	}

	return ts[0], nil
}

func (t *Core[T]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	var ts []T

	cursor, err := t.c.Find(ctx, filter, opts...)
	if err != nil {
		return ts, err
	}

	err = cursor.All(ctx, &ts)
	if err != nil {
		return ts, err
	}

	return ts, nil
}

func (t *Core[T]) FindOne(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (T, error) {
	var zero T

	var ts []T

	cursor, err := t.c.Find(ctx, filter, opts...)
	if err != nil {
		return zero, err
	}

	err = cursor.All(ctx, &ts)
	if err != nil {
		return zero, err
	}

	if len(ts) == 0 {
		return zero, nil
	}

	return ts[0], nil
}

func (t *Core[T]) InsertNX(ctx context.Context, doc T, filter bson.M) (string, bool, error) {

	list, err := t.List(ctx, filter)
	if err != nil {
		return "", false, err
	}

	if len(list) != 0 {
		return list[0].GetID(), false, nil
	}

	res, err := t.c.InsertOne(ctx, doc)
	if err != nil {
		return "", false, err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), true, nil

}

func (t *Core[T]) ReplaceByID(ctx context.Context, id string, doc T) (bool, error) {

	doc.SetID("")

	res, err := t.c.ReplaceOne(ctx, bson.M{"_id": ObjectId(id)}, doc)
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

// Deprecated: UpdateByIDXX
func (t *Core[T]) UpdateByIDXX(ctx context.Context, id string, update bson.M) (bool, error) {

	one, err := t.c.UpdateOne(ctx, bson.M{"_id": ObjectId(id)}, update, options.Update().SetUpsert(false))

	//one, err := t.c.UpdateByID(ctx, bson.M{"_id": ObjectId(id)}, update, options.Update().SetUpsert(false))
	if err != nil {
		return false, err
	}

	return one.ModifiedCount > 0, nil
}
func (t *Core[T]) UpdateByIDIfExists(ctx context.Context, id string, update *UpdateOperation) (bool, error) {

	u := update.ToBson()

	one, err := t.c.UpdateOne(ctx, bson.M{"_id": ObjectId(id)}, u, options.Update().SetUpsert(false))

	//one, err := t.c.UpdateByID(ctx, bson.M{"_id": ObjectId(id)}, update, options.Update().SetUpsert(false))
	if err != nil {
		return false, err
	}

	return one.ModifiedCount > 0, nil
}
func (t *Core[T]) UpdateByIDsIfExists(ctx context.Context, ids []string, update *UpdateOperation) (bool, error) {

	u := update.ToBson()

	one, err := t.c.UpdateMany(ctx, Filter().Ids(ids).B(), u, options.Update().SetUpsert(false))

	//one, err := t.c.UpdateByID(ctx, bson.M{"_id": ObjectId(id)}, update, options.Update().SetUpsert(false))
	if err != nil {
		return false, err
	}

	return one.ModifiedCount > 0, nil
}

func (t *Core[T]) UpdateOne(ctx context.Context, filter bson.M, update *UpdateOperation) (bool, error) {

	u := update.ToBson()

	one, err := t.c.UpdateOne(ctx, filter, u, options.Update().SetUpsert(false))

	//one, err := t.c.UpdateByID(ctx, bson.M{"_id": ObjectId(id)}, update, options.Update().SetUpsert(false))
	if err != nil {
		return false, err
	}

	return one.ModifiedCount > 0, nil
}

func (t *Core[T]) UpdateByIDsXX(ctx context.Context, ids []string, update bson.M) (bool, error) {

	one, err := t.c.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ObjectIds(ids)}}, update, options.Update().SetUpsert(false))

	//one, err := t.c.UpdateByID(ctx, bson.M{"_id": ObjectId(id)}, update, options.Update().SetUpsert(false))
	if err != nil {
		return false, err
	}

	return one.ModifiedCount > 0, nil
}

func (t *Core[T]) InsertMany(ctx context.Context, docs ...T) ([]string, error) {

	var ds []interface{}
	for _, doc := range docs {
		ds = append(ds, doc)
	}

	if len(ds) == 0 {
		return []string{}, nil
	}

	res, err := t.c.InsertMany(ctx, ds)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, id := range res.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}

	return ids, nil
}

func (t *Core[T]) Insert(ctx context.Context, doc T) (T, error) {

	var zero T

	doc.SetID("")

	one, err := t.c.InsertOne(ctx, doc)
	if err != nil {
		return zero, err
	}

	doc.SetID(one.InsertedID.(primitive.ObjectID).Hex())

	return doc, nil
}

func (t *Core[T]) List(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error) {

	//opt := options.Find()
	//if len(opts) > 0 {
	//	opt = opts[0]
	//}

	filter["deletedAt"] = bson.M{"$exists": false}

	cursor, err := t.c.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var items []T

	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil

}

func (t *Core[T]) ListAndCount(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, int64, error) {

	list, err := t.List(ctx, filter, opts...)
	if err != nil {
		return nil, 0, err
	}

	count, err := t.c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil

}

type DeleteOptions struct {
	Reserve bool
}

func (t *Core[T]) DeleteByID(ctx context.Context, id string, options ...DeleteOptions) error {

	_, err := t.Delete(ctx, bson.M{"_id": ObjectId(id)}, options...)

	if err != nil {
		return err
	}

	return nil
}

func (t *Core[T]) Delete(ctx context.Context, filter bson.M, options ...DeleteOptions) (int64, error) {

	if len(options) > 0 && options[0].Reserve {
		many, err := t.c.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"deletedAt": time.Now().Unix()}})
		if err != nil {
			return 0, err
		}

		return many.ModifiedCount, nil

	} else {

		d, err := t.c.DeleteMany(ctx, filter)
		if err != nil {
			return 0, err
		}
		return d.DeletedCount, nil
	}
}

func (t *Core[T]) Count(ctx context.Context, filter bson.M) (int64, error) {

	d, err := t.c.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return d, nil
}

func (t *Core[T]) C() *mongo.Collection {
	return t.c
}
func (t *Core[T]) D() *mongo.Database {
	return t.db
}

// ListWithFilterAndSort returns items matching the filter with sorting and pagination
func (t *Core[T]) ListWithFilterAndSort(ctx context.Context, filter bson.M, sort bson.M, page, size int64) ([]T, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	filter["deletedAt"] = bson.M{"$exists": false}

	// Count total
	total, err := t.c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build find options
	findOpts := options.Find().
		SetSort(sort).
		SetSkip((page - 1) * size).
		SetLimit(size)

	cursor, err := t.c.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}

	var items []T
	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

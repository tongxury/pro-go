package mgz

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CoreCollectionV2[T any] struct {
	c *mongo.Collection
}

func NewCoreCollectionV2[T any](db *mongo.Database, name string) *CoreCollectionV2[T] {
	return &CoreCollectionV2[T]{
		c: db.Collection(name, options.Collection().
			SetBSONOptions(&options.BSONOptions{UseJSONStructTags: true})),
	}
}

func (t *CoreCollectionV2[T]) ReplaceOneXX(ctx context.Context, filter bson.M, document *T) error {
	return t.c.FindOneAndReplace(ctx, filter, document).Err()
}

func (t *CoreCollectionV2[T]) UpdateOneXX(ctx context.Context, filter bson.M, fields bson.M) (bool, error) {
	rsp, err := t.c.UpdateOne(ctx, filter, bson.M{"$set": fields})
	if err != nil {
		return false, err
	}

	return rsp.ModifiedCount > 0 || rsp.UpsertedCount > 0, nil
}

func (t *CoreCollectionV2[T]) FindAndUpdateOne(ctx context.Context, filter bson.M, fields bson.M) (*T, error) {

	var result T

	err := t.c.FindOneAndUpdate(ctx,
		filter,
		bson.D{{"$set", fields}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *CoreCollectionV2[T]) InsertNX(ctx context.Context, filter bson.M, item *T) (bool, error) {

	result, err := t.c.UpdateOne(ctx, filter, bson.M{"$setOnInsert": item},
		options.Update().SetUpsert(true))

	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, nil
}

func (t *CoreCollectionV2[T]) Insert(ctx context.Context, document *T) (*T, error) {

	//// 使用反射设置 ID
	//val := reflect.ValueOf(document).Elem()
	//idField := val.FieldByName("ID")
	//
	//if !idField.IsValid() {
	//	return nil, errors.New("ID field not found")
	//}
	//
	//if !idField.CanSet() {
	//	return nil, errors.New("ID field cannot be set")
	//}
	//
	//// 生成新的 ObjectID
	//newID := primitive.NewObjectID()
	//idField.Set(reflect.ValueOf(newID))

	_, err := t.c.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (t *CoreCollectionV2[T]) List(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, error) {

	if len(filter) == 0 {
		filter = bson.M{}
	}
	filter["deletedAt"] = bson.D{{"$exists", false}}

	cursor, err := t.c.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var items []*T

	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil

}

func (t *CoreCollectionV2[T]) ListAndCount(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, int64, error) {

	filter["deletedAt"] = bson.D{{"$exists", false}}

	//defaultSort := bson.B{{"$sort", bson.B{{"createdAt", -1}}}}
	//
	//defaultOptions := &options.FindOptions{
	//	Sort: defaultSort,
	//}
	//if len(opts) > 0 {
	//	defaultOptions = opts[0]
	//}

	opt := options.Find()
	if len(opts) > 0 {
		opt = opts[0]
	}

	//opt.Sort = sort

	cursor, err := t.c.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}

	var items []*T

	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, 0, err
	}

	count, err := t.c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil

}

func (t *CoreCollectionV2[T]) Destroy(ctx context.Context, filter bson.M) (int64, error) {

	res, err := t.c.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, err
}

package mgz

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"store/pkg/sdk/conv"
	"time"
)

type CoreCollection[T any] struct {
	*mongo.Collection
}

func NewCoreCollection[T any](db *mongo.Database, name string) *CoreCollection[T] {
	return &CoreCollection[T]{
		Collection: db.Collection(name,
			options.Collection().
				SetBSONOptions(&options.BSONOptions{UseJSONStructTags: true}),
		),
	}
}

func (t *CoreCollection[T]) InsertBatchSkipFailed(ctx context.Context, documents []*T) error {
	if len(documents) == 0 {
		return nil
	}
	_, err := t.InsertMany(ctx, conv.AnySlice(documents), options.InsertMany().SetOrdered(false))
	return err
}
func (t *CoreCollection[T]) Insert(ctx context.Context, document *T) error {
	if document == nil {
		return nil
	}
	_, err := t.InsertOne(ctx, document)
	return err
}

func (t *CoreCollection[T]) InsertBatch(ctx context.Context, documents []*T) error {
	if len(documents) == 0 {
		return nil
	}
	_, err := t.InsertMany(ctx, conv.AnySlice(documents))
	return err
}

func (t *CoreCollection[T]) InsertFieldIfNotExists(ctx context.Context, filters bson.M, fieldKey string, value any) (bool, error) {

	filters[fieldKey] = bson.M{"$exists": false}

	result, err := t.UpdateOne(ctx,
		filters,
		bson.M{
			"$set": bson.M{
				fieldKey: value,
			},
		},
	)

	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

func (t *CoreCollection[T]) ReplaceOrInsertOneByID(ctx context.Context, id string, doc *T) (bool, error) {
	if doc == nil {
		return false, nil
	}
	return t.ReplaceOrInsertOne(ctx, bson.M{"_id": id}, doc)
}

func (t *CoreCollection[T]) ReplaceOrInsertOne(ctx context.Context, filter bson.M, doc *T) (bool, error) {
	if doc == nil {
		return false, nil
	}
	result, err := t.ReplaceOne(ctx, filter, doc,
		options.Replace().SetUpsert(true))

	if err != nil {
		return false, err
	}

	return result.MatchedCount > 0, nil
}

func (t *CoreCollection[T]) GetById(ctx context.Context, id string) (*T, error) {
	user, err := t.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf(` %s not found`, id)
	}

	return user, nil
}

func (t *CoreCollection[T]) FindById(ctx context.Context, id string) (*T, error) {

	var ts []*T

	cursor, err := t.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &ts)
	if err != nil {
		return nil, err
	}

	if len(ts) == 0 {
		return nil, nil
	}

	return ts[0], nil
}

func (t *CoreCollection[T]) FindByIds(ctx context.Context, ids []string) ([]*T, error) {

	var ts []*T

	cursor, err := t.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &ts)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (t *CoreCollection[T]) InsertIfNotExistsByID(ctx context.Context, id string, item *T) (bool, error) {
	return t.InsertIfNotExists(ctx, bson.M{"_id": id}, item)
}

func (t *CoreCollection[T]) InsertIfNotExists(ctx context.Context, filter bson.M, item *T) (bool, error) {

	result, err := t.UpdateOne(ctx, filter, bson.M{"$setOnInsert": item},
		options.Update().SetUpsert(true))

	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, nil
}

func (t *CoreCollection[T]) Archive(ctx context.Context, filter bson.M) error {

	_, err := t.UpdateOne(ctx, filter,
		bson.D{{"$set", bson.D{{"deletedAt", time.Now().Unix()}}}})
	if err != nil {
		return err
	}

	return nil
}

func (t *CoreCollection[T]) ArchiveByID(ctx context.Context, id string) error {

	_, err := t.UpdateByID(ctx, id,
		bson.D{{"$set", bson.D{{"deletedAt", time.Now().Unix()}}}})
	if err != nil {
		return err
	}

	return nil
}

func (t *CoreCollection[T]) DestroyById(ctx context.Context, id string) error {

	filter := bson.D{{"_id", id}}
	_, err := t.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (t *CoreCollection[T]) UpdateFieldsById(ctx context.Context, id string, update bson.M) (bool, error) {
	return t.UpdateFields(ctx, bson.M{"_id": id}, update)
}

func (t *CoreCollection[T]) UpdateFields(ctx context.Context, filter bson.M, update bson.M) (bool, error) {
	rsp, err := t.UpdateMany(ctx, filter, bson.M{"$set": update},
		options.Update().SetUpsert(false))
	if err != nil {
		return false, err
	}

	return rsp.ModifiedCount > 0 || rsp.UpsertedCount > 0, nil
}

func (t *CoreCollection[T]) UpdateOneById(ctx context.Context, id string, key string, value any) (bool, error) {

	rsp, err := t.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{key: value}})
	if err != nil {
		return false, err
	}

	return rsp.UpsertedCount > 0, nil
}

func (t *CoreCollection[T]) Count(ctx context.Context, filter bson.M) (int64, error) {
	if len(filter) == 0 {
		filter = bson.M{}
	}
	filter["deletedAt"] = bson.D{{"$exists", false}}

	return t.CountDocuments(ctx, filter)
}

//func (t *CoreCollection[T]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, error) {
//	return t.List(ctx, filter, opts...)
//}

func (t *CoreCollection[T]) List(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, error) {

	if len(filter) == 0 {
		filter = bson.M{}
	}
	filter["deletedAt"] = bson.D{{"$exists", false}}

	cursor, err := t.Find(ctx, filter, opts...)
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

func (t *CoreCollection[T]) ListAndCount(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, int64, error) {

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

	cursor, err := t.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}

	var items []*T

	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, 0, err
	}

	count, err := t.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil

}

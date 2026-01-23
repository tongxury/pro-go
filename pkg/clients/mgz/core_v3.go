package mgz

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CoreCollectionV3[T any] struct {
	c  *mongo.Collection
	db *mongo.Database
}

func NewCoreCollectionV3[T any](db *mongo.Database, name string) *CoreCollectionV3[T] {
	return &CoreCollectionV3[T]{
		db: db,
		c: db.Collection(name, options.Collection().
			SetBSONOptions(&options.BSONOptions{UseJSONStructTags: true})),
	}
}

func (t *CoreCollectionV3[T]) GetCoreCollection() *mongo.Collection {
	return t.c
}
func (t *CoreCollectionV3[T]) GetDatabase() *mongo.Database {
	return t.db
}

func (t *CoreCollectionV3[T]) ReplaceOneXX(ctx context.Context, filter bson.M, document *T) error {
	return t.c.FindOneAndReplace(ctx, filter, document).Err()
}

func (t *CoreCollectionV3[T]) UpdateOneXX(ctx context.Context, filter bson.M, fields bson.M) (bool, error) {
	rsp, err := t.c.UpdateOne(ctx, filter, bson.M{"$set": fields})
	if err != nil {
		return false, err
	}

	return rsp.ModifiedCount > 0 || rsp.UpsertedCount > 0, nil
}

func (t *CoreCollectionV3[T]) UpdateOneXXPush(ctx context.Context, filter bson.M, field string, newValues []bson.M, position int) (bool, error) {

	update := bson.M{
		"$push": bson.M{
			field: bson.M{
				"$each":     newValues,
				"$position": position, // 在数组开头插入
			},
		},
	}

	rsp, err := t.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return rsp.ModifiedCount > 0 || rsp.UpsertedCount > 0, nil
}

func (t *CoreCollectionV3[T]) UpdateOneXXPushById(ctx context.Context, id string, field string, newValues []bson.M, position int) (bool, error) {
	return t.UpdateOneXXPush(ctx, bson.M{"_id": ObjectId(id)}, field, newValues, position)
}

func (t *CoreCollectionV3[T]) UpdateOneXXById(ctx context.Context, id string, fields bson.M) (bool, error) {
	rsp, err := t.c.UpdateOne(ctx, bson.M{"_id": ObjectId(id)}, bson.M{"$set": fields})
	if err != nil {
		return false, err
	}

	return rsp.ModifiedCount > 0 || rsp.UpsertedCount > 0, nil
}

func (t *CoreCollectionV3[T]) FindAndUpdateOne(ctx context.Context, filter bson.M, fields bson.M) (*T, error) {

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

//
//func (t *CoreCollectionV3[T]) InsertInBatch(ctx context.Context, documents []*T) ([]string, error)  {
//
//	var docs []interface{}
//	for _, x := range documents {
//		docs = append(docs, x)
//	}
//
//	many, err := t.c.InsertMany(ctx, docs)
//	if err != nil {
//		return nil, err
//	}
//
//	return many.InsertedIDs, nil
//}

func (t *CoreCollectionV3[T]) InsertNX(ctx context.Context, filter bson.M, item *T) (bool, string, error) {

	//t.c.Database().Client().StartSession()

	result, err := t.c.UpdateOne(ctx, filter, bson.M{"$setOnInsert": item},
		options.Update().SetUpsert(true))

	if err != nil {
		return false, "", err
	}

	if result.UpsertedCount > 0 {
		return result.UpsertedCount > 0, result.UpsertedID.(primitive.ObjectID).Hex(), nil
	}

	return false, "", nil
}

func (t *CoreCollectionV3[T]) GetById(ctx context.Context, id string) (*T, error) {
	user, err := t.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf(` %s not found`, id)
	}

	return user, nil
}

func (t *CoreCollectionV3[T]) FindById(ctx context.Context, id string) (*T, error) {

	var ts []*T

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cursor, err := t.c.Find(ctx, bson.M{"_id": hex})
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

func (t *CoreCollectionV3[T]) Insert(ctx context.Context, document *T) (string, error) {

	id, err := t.c.InsertOne(ctx, document)

	if err != nil {
		return "", err
	}

	return id.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (t *CoreCollectionV3[T]) Inserts(ctx context.Context, document []*T) ([]string, error) {

	ids, err := t.c.InsertMany(ctx, conv.AnySlice(document))

	if err != nil {
		return nil, err
	}

	var result []string
	for _, x := range ids.InsertedIDs {
		result = append(result, x.(primitive.ObjectID).Hex())
	}

	return result, nil
}

func (t *CoreCollectionV3[T]) List(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, error) {

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

func (t *CoreCollectionV3[T]) ListAndCount(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*T, int64, error) {

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

func (t *CoreCollectionV3[T]) Destroy(ctx context.Context, filter bson.M) (int64, error) {

	res, err := t.c.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, err
}

func (t *CoreCollectionV3[T]) DestroyById(ctx context.Context, id string) (int64, error) {

	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	res, err := t.c.DeleteMany(ctx, bson.M{"_id": hexId})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, err
}

func (t *CoreCollectionV3[T]) Count(ctx context.Context, filter bson.M) (int64, error) {
	return t.c.CountDocuments(ctx, filter)
}

func (t *CoreCollectionV3[T]) Sum(ctx context.Context, filter bson.M, field string) (float64, error) {
	// 添加软删除过滤
	if len(filter) == 0 {
		filter = bson.M{}
	}
	filter["deletedAt"] = bson.D{{"$exists", false}}

	// 构建聚合管道
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$" + field},
		}},
	}

	cursor, err := t.c.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	// 如果没有结果，返回0
	if len(results) == 0 {
		return 0, nil
	}

	// 获取求和结果
	total, ok := results[0]["total"]
	if !ok {
		return 0, nil
	}

	// 根据字段类型转换结果
	switch v := total.(type) {
	case float64:
		return v, nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported field type for sum: %T", total)
	}
}

// 如果需要支持多个字段求和，可以使用这个方法
func (t *CoreCollectionV3[T]) SumMultiple(ctx context.Context, filter bson.M, fields []string) (map[string]float64, error) {
	// 添加软删除过滤
	if len(filter) == 0 {
		filter = bson.M{}
	}
	filter["deletedAt"] = bson.D{{"$exists", false}}

	// 构建聚合管道
	groupFields := bson.M{"_id": nil}
	for _, field := range fields {
		groupFields["sum_"+field] = bson.M{"$sum": "$" + field}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$group": groupFields},
	}

	cursor, err := t.c.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// 如果没有结果，返回空map
	if len(results) == 0 {
		result := make(map[string]float64)
		for _, field := range fields {
			result[field] = 0
		}
		return result, nil
	}

	// 解析结果
	result := make(map[string]float64)
	for _, field := range fields {
		sumKey := "sum_" + field
		if total, ok := results[0][sumKey]; ok {
			switch v := total.(type) {
			case float64:
				result[field] = v
			case int32:
				result[field] = float64(v)
			case int64:
				result[field] = float64(v)
			case int:
				result[field] = float64(v)
			default:
				result[field] = 0
			}
		} else {
			result[field] = 0
		}
	}

	return result, nil
}

package mgz

import (
	"store/pkg/sdk/helper/mathz"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ToMap[T IColl](src []T) map[string]T {

	mp := make(map[string]T, len(src))
	for _, t := range src {
		mp[t.GetID()] = t
	}

	return mp
}

func Ids[T IColl](src []T) []string {

	var ids []string
	for _, t := range src {
		if t.GetID() == "" {
			continue
		}
		ids = append(ids, t.GetID())
	}

	return ids
}

func Paging[T mathz.Int](page, size T) *options.FindOptions {

	opts := options.Find()

	if size == 0 {
		return opts
	}

	if page == 0 {
		return opts.SetLimit(int64(size))
	}

	return opts.SetLimit(int64(size)).SetSkip(int64((page - 1) * size))
}

func ObjectId(id string) primitive.ObjectID {

	hex, _ := primitive.ObjectIDFromHex(id)

	return hex
}

func ObjectIds(ids []string) []primitive.ObjectID {

	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectIds = append(objectIds, ObjectId(id))
	}

	return objectIds
}

func Set(fields bson.M) bson.M {
	return bson.M{"$set": fields}
}

func Exists(field string) bson.M {
	return bson.M{field: bson.M{"$exists": true}}
}

func NotExists(field string) bson.M {
	return bson.M{field: bson.M{"$exists": false}}
}

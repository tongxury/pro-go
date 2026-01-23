package redizv2

type RedisJSON[T Document] struct {
	*FTSearch[T]
	*RJSON[T]
}

func NewRedisJSON[T Document](client *RedisClient, keyPrefix string) *RedisJSON[T] {
	return &RedisJSON[T]{
		FTSearch: NewFTSearch[T](client),
		RJSON:    NewRJSON[T](client, keyPrefix),
	}
}

//func (r RedisJsonDB[T]) InsertMany(ctx context.Context, values ...T) error {
//
//	if len(values) == 0 {
//		return nil
//	}
//
//	args := helper.Mapping(values, func(x T) redis.JSONSetArgs {
//
//		return redis.JSONSetArgs{
//			Key:   r.keyPrefix + x.GetKey(),
//			Path:  "$",
//			Value: x,
//		}
//	})
//
//	_, err := r.c.JSONMSetArgs(ctx, args).Result()
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}

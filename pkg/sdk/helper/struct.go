package helper

import (
	"reflect"
)

// MergeStructs 合并多个相同类型的对象，优先取后面传入的对象的非空值
func MergeStructs[T any](objs ...*T) *T {
	if len(objs) == 0 {
		var zero *T
		return zero
	}

	// 获取第一个对象的类型，并检查是否为 nil
	objType := reflect.TypeOf(objs[0])
	if objType.Kind() != reflect.Pointer || objType.Elem().Kind() != reflect.Struct {
		//panic("All inputs must be pointers to structs")
		return nil
	}

	// 创建一个空的结构体实例的指针
	merged := reflect.New(objType.Elem())

	// 从最后一个对象开始，向前迭代
	for i := len(objs) - 1; i >= 0; i-- {
		if objs[i] == nil {
			continue // 忽略 nil 参数
		}
		val := reflect.ValueOf(objs[i]).Elem()

		// 遍历字段
		for j := 0; j < val.NumField(); j++ {
			field := val.Field(j)
			// 只有当 merged 的字段是零值时，才从当前对象复制字段
			if merged.Elem().Field(j).IsZero() {
				merged.Elem().Field(j).Set(field)
			}
		}
	}

	return merged.Interface().(*T) // 返回合并后的对象指针
}

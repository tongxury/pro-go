package mgz

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type UpdateOperation struct {
	set   bson.M
	unset bson.M
	inc   bson.M
	push  bson.M
	pull  bson.M
}

// Op 创建一个新的 UpdateOperation
func Op() *UpdateOperation {
	return &UpdateOperation{
		set:   make(bson.M),
		unset: make(bson.M),
		inc:   make(bson.M),
		push:  make(bson.M),
		pull:  make(bson.M),
	}
}

// SetField 设置字段值
func (u *UpdateOperation) Set(key string, value interface{}) *UpdateOperation {
	if u.set == nil {
		u.set = make(bson.M)
	}
	u.set[key] = value
	return u
}

// SetField 设置字段值
func (u *UpdateOperation) SetListItem(key string, index int, field string, value interface{}) *UpdateOperation {
	if u.set == nil {
		u.set = make(bson.M)
	}
	u.set[fmt.Sprintf("%s.%d.%s", key, index, field)] = value
	return u
}

// UnsetField 删除字段
func (u *UpdateOperation) Unset(key string) *UpdateOperation {
	if u.unset == nil {
		u.unset = make(bson.M)
	}
	u.unset[key] = ""
	return u
}

// IncField 增加字段值（支持数字类型）
func (u *UpdateOperation) Inc(key string, value interface{}) *UpdateOperation {
	if u.inc == nil {
		u.inc = make(bson.M)
	}
	u.inc[key] = value
	return u
}

// PushField 向数组字段添加元素
func (u *UpdateOperation) Push(key string, value interface{}) *UpdateOperation {
	if u.push == nil {
		u.push = make(bson.M)
	}
	u.push[key] = value
	return u
}

// PullField 从数组字段移除元素
func (u *UpdateOperation) Pull(key string, value interface{}) *UpdateOperation {
	if u.pull == nil {
		u.pull = make(bson.M)
	}
	u.pull[key] = value
	return u
}

// ToBson 转换为 bson.M
func (u *UpdateOperation) ToBson() bson.M {
	result := bson.M{}

	if len(u.set) > 0 {
		result["$set"] = u.set
	}
	if len(u.unset) > 0 {
		result["$unset"] = u.unset
	}
	if len(u.inc) > 0 {
		result["$inc"] = u.inc
	}
	if len(u.push) > 0 {
		result["$push"] = u.push
	}
	if len(u.pull) > 0 {
		result["$pull"] = u.pull
	}

	return result
}

// IsEmpty 检查是否有任何更新操作
func (u *UpdateOperation) IsEmpty() bool {
	return len(u.set) == 0 && len(u.unset) == 0 && len(u.inc) == 0 &&
		len(u.push) == 0 && len(u.pull) == 0
}

// SetFields 批量设置字段
func (u *UpdateOperation) Sets(fields bson.M) *UpdateOperation {
	if u.set == nil {
		u.set = make(bson.M)
	}
	for k, v := range fields {
		u.set[k] = v
	}
	return u
}

// SetIf 条件设置
func (u *UpdateOperation) SetIf(condition bool, key string, value interface{}) *UpdateOperation {
	if condition {
		u.Set(key, value)
	}
	return u
}

// SetIfNotNil 如果值不为 nil 则设置
func (u *UpdateOperation) SetIfNotNil(key string, value interface{}) *UpdateOperation {
	if value != nil {
		u.Set(key, value)
	}
	return u
}

// PushAll 向数组添加多个元素
func (u *UpdateOperation) PushAll(key string, values ...interface{}) *UpdateOperation {
	if u.push == nil {
		u.push = make(bson.M)
	}
	u.push[key] = bson.M{"$each": values}
	return u
}

// PullAll 从数组移除多个元素
func (u *UpdateOperation) PullAll(key string, values ...interface{}) *UpdateOperation {
	if u.pull == nil {
		u.pull = make(bson.M)
	}
	u.pull[key] = bson.M{"$in": values}
	return u
}

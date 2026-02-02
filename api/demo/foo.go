package demopb

// extra.go 文件用于扩展 Protobuf 自动生成的 Go 结构体。
// 这里的包名需要对应 proto 文件中 go_package 选项分号后的包名。

// --- Foo 扩展 ---

// GetID 是一个辅助方法，用于处理 ID 映射。
func (m *Foo) GetID() string {
	if m == nil {
		return ""
	}
	return m.XId
}

// SetID 设置 ID。
func (m *Foo) SetID(id string) {
	if m == nil {
		return
	}
	m.XId = id
}

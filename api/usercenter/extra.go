package ucpb

func (m *User) GetID() string {
	return m.XId
}

func (m *User) SetID(id string) {
	m.XId = id
}

package voiceagent

func (x *Assessment) GetID() string {
	if x != nil {
		return x.XId
	}
	return ""
}

func (x *Assessment) SetID(id string) {
	if x != nil {
		x.XId = id
	}
}

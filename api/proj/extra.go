package projpb

func (x *Resource) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *Resource) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (x *ResourceSegment) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *ResourceSegment) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (x *Commodity) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *Commodity) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (x *Asset) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *Asset) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (m *AppSettings) GetID() string {
	if m == nil {
		return ""
	}
	return m.XId
}

func (m *AppSettings) SetID(id string) {
	if m == nil {
		return
	}
	m.XId = id
}

func (m *Task) GetID() string {
	if m == nil {
		return ""
	}
	return m.XId
}

func (m *Task) SetID(id string) {
	if m == nil {
		return
	}
	m.XId = id
}

func (x *Feedback) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *Feedback) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (x *Asset) EnsureAttrs() *Asset_Attrs {

	attrs := &Asset_Attrs{}

	if x == nil {
		return attrs
	}
	if x.Attrs == nil {
		x.Attrs = attrs
	}
	return x.Attrs
}

func (x *Session) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *Session) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (x *SessionSegment) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *SessionSegment) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (x *Workflow) GetID() string {
	if x == nil {
		return ""
	}
	return x.XId
}

func (x *Workflow) SetID(id string) {
	if x == nil {
		return
	}
	x.XId = id
}

func (t *TypedTags) Tags() []string {

	if t == nil {
		return []string{}
	}

	var tags []string
	tags = append(tags, t.GetAction()...)
	tags = append(tags, t.GetEmotion()...)
	tags = append(tags, t.GetPerson()...)
	tags = append(tags, t.GetFocusOn()...)
	tags = append(tags, t.GetPicture()...)
	tags = append(tags, t.GetScene()...)
	tags = append(tags, t.GetText()...)
	tags = append(tags, t.GetShootingStyle()...)

	return tags
}

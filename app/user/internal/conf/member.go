package conf

type Description struct {
	Name  string `json:"name,omitempty"`
	Desc  string `json:"desc,omitempty"`
	Label string `json:"label,omitempty"`
}

type Descriptions []*Description

func (ts Descriptions) ByName(name string) *Description {

	for _, t := range ts {
		if t.Name == name {
			return t
		}
	}

	return nil
}

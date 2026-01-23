package repo

type Promotion struct {
	Id       string
	Limit    int64
	Supports Supports
}

type Supports []*Support

type Support struct {
	Level string
	Cycle string
}

type Promotions []*Promotion

func (ts Promotions) AsMap() map[string]*Promotion {

	rsp := make(map[string]*Promotion, len(ts))

	for _, t := range ts {
		rsp[t.Id] = t
	}

	return rsp
}

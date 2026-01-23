package dexpb

type Users []*User

func (ts Users) AsMap() map[string]*User {

	var mp = make(map[string]*User)
	for _, x := range ts {
		mp[x.XId] = x
	}

	return mp
}

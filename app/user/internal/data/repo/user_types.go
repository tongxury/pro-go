package repo

import (
	"store/app/user/internal/data/repo/ent"
)

type User struct {
	ent.User
}

type Users []*User

func (ts Users) IDs() []int64 {
	var rsp []int64
	for _, t := range ts {
		rsp = append(rsp, t.ID)
	}
	return nil
}

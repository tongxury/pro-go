package repo

import (
	"time"
)

type User struct {
	ID          int64     `ch:"id"`
	Nickname    string    `ch:"nickname"`
	Email       string    `ch:"email"`
	CreatedAt   time.Time `ch:"create_time"`
	LastLoginAt time.Time `ch:"last_login_time"`
	Password    string    `ch:"pass_word"`
	//Verify      int8      `ck:"verify"`
}

type Users []User

func (ts Users) IDs() []int64 {
	var ids []int64
	for _, t := range ts {
		ids = append(ids, t.ID)
	}

	return ids
}

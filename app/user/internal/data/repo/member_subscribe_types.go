package repo

import "store/app/user/internal/data/repo/ent"

type MemberSubscribe struct {
	*ent.MemberSubscribe
}

type MemberSubscribes []*MemberSubscribe

func (ts MemberSubscribes) GroupByUserID() map[int64]MemberSubscribes {

	rsp := make(map[int64]MemberSubscribes, len(ts))

	for _, t := range ts {
		rsp[t.UserID] = append(rsp[t.UserID], t)
	}

	return rsp
}

func AsMemberSubscribes(subs []*ent.MemberSubscribe) MemberSubscribes {
	var rsp MemberSubscribes
	for _, sub := range subs {
		rsp = append(rsp, &MemberSubscribe{
			MemberSubscribe: sub,
		})
	}

	return rsp
}

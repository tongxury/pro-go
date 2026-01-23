package repo

type State struct {
	Channel string `ch:"channel"`
	Count   uint64 `ch:"count"`
}

type States []*State

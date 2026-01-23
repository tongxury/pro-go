package repo

type UserFunnelRecord struct {
	Date    string `ch:"date"`
	Channel string `ch:"channel"`
	Event   string `ch:"event"`
	Level   uint8  `ch:"level"`
	Count   uint64 `ch:"count"`
}

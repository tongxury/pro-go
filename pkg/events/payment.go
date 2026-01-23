package events

type PaymentSuccessEvent struct {
	UniqueId string
	UserID   string
	DeviceID string
	Amount   float64
	Platform string
	Ts       int64
}

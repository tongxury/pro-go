package enums

type SubscribeEvent string

const (
	SubscribeEvent_Started          PaymentStatus = "started"
	SubscribeEvent_ExpiredAtUpdated PaymentStatus = "expiredAtUpdated"
	SubscribeEvent_Canceled         PaymentStatus = "canceled"
)

func (t SubscribeEvent) Values() []string {
	return []string{
		SubscribeEvent_Started.String(),
		SubscribeEvent_ExpiredAtUpdated.String(),
		SubscribeEvent_Canceled.String(),
	}
}

func (t SubscribeEvent) String() string {
	return string(t)
}

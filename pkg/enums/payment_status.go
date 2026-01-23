package enums

type PaymentStatus string

const (
	PaymentStatus_Created   PaymentStatus = "created"
	PaymentStatus_Expired   PaymentStatus = "expired"
	PaymentStatus_Complete  PaymentStatus = "complete"
	PaymentStatus_Failed    PaymentStatus = "failed"
	PaymentStatus_Cancelled PaymentStatus = "cancelled"
)

func (t PaymentStatus) Values() []string {
	return []string{
		PaymentStatus_Created.String(),
		PaymentStatus_Expired.String(),
		PaymentStatus_Complete.String(),
		PaymentStatus_Failed.String(),
		PaymentStatus_Cancelled.String(),
	}
}

func (t PaymentStatus) String() string {
	return string(t)
}

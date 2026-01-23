package enums

type MemberSubscribeStatus string

const (
	MemberSubscribeStatus_Pending     MemberSubscribeStatus = "pending"
	MemberSubscribeStatus_Subscribing MemberSubscribeStatus = "subscribing"
	MemberSubscribeStatus_Unsubscribe MemberSubscribeStatus = "unsubscribe"
	MemberSubscribeStatus_Expired     MemberSubscribeStatus = "expired"
)

func (MemberSubscribeStatus) Values() []string {
	return []string{
		MemberSubscribeStatus_Pending.String(),
		MemberSubscribeStatus_Subscribing.String(),
		MemberSubscribeStatus_Unsubscribe.String(),
		MemberSubscribeStatus_Expired.String(),
	}
}

func (t MemberSubscribeStatus) String() string {
	return string(t)
}

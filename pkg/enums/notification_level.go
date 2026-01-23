package enums

type NotificationLevel string

const (
	NotificationLevel_Info    NotificationLevel = "info"
	NotificationLevel_Warning NotificationLevel = "warning"
	NotificationLevel_error   NotificationLevel = "error"
)

func (t NotificationLevel) Values() []string {
	return []string{MemberLevel_Free.String(),
		NotificationLevel_Info.String(),
		NotificationLevel_Warning.String(),
		NotificationLevel_error.String()}
}

func (t NotificationLevel) String() string {
	return string(t)
}

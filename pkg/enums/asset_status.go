package enums

type AssetStatus string

const (
	AssetStatus_Pending AssetStatus = "pending"
	AssetStatus_Normal  AssetStatus = "normal"
	AssetStatus_Discard AssetStatus = "discard"
)

func (t AssetStatus) Values() []string {
	return []string{
		AssetStatus_Pending.String(),
		AssetStatus_Normal.String(),
		AssetStatus_Discard.String(),
	}
}

func (t AssetStatus) String() string {
	return string(t)
}

func (t AssetStatus) Name() string {
	switch t {
	case AssetStatus_Pending:
		return "In Verification"
	case AssetStatus_Normal:
		return "Approved"
	case AssetStatus_Discard:
		return "Discard"
	default:
		return ""
	}
}

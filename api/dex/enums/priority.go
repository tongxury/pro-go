package enums

var (
	PriorityLevel_Medium   = "medium"
	PriorityLevel_High     = "high"
	PriorityLevel_VeryHigh = "veryHigh"

	PriorityLevels = []string{PriorityLevel_Medium, PriorityLevel_High, PriorityLevel_VeryHigh}
)

//
//var (
//	Priority_Medium = &tgbotpb.OrderSettings_Priority{
//		Value: PriorityLevel_Medium,
//		//Max:   0.006,
//	}
//	Priority_High = &tgbotpb.OrderSettings_Priority{
//		Value: PriorityLevel_High,
//		//Max:   0.01,
//	}
//	Priority_VeryHigh = &tgbotpb.OrderSettings_Priority{
//		Value: PriorityLevel_VeryHigh,
//		//Max:   0.015,
//	}
//)
//
//func PriorityByName(value string) *tgbotpb.OrderSettings_Priority {
//
//	switch value {
//	case PriorityLevel_Medium:
//		return Priority_Medium
//	case PriorityLevel_High:
//		return Priority_High
//	case PriorityLevel_VeryHigh:
//		return Priority_VeryHigh
//	default:
//		return nil
//	}
//}

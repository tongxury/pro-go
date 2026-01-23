package configs

type Plan struct {
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	Amount         float64 `json:"amount"`
	OriginAmount   float64 `json:"originAmount,omitempty"`
	Months         int     `json:"months"`
	CreditPerMonth int64   `json:"creditPerMonth,omitempty"`
}

var plans = []Plan{
	{Id: "single_year_unlimited", Title: "12个月(无限)", Amount: 938, OriginAmount: 1176, Months: 12, CreditPerMonth: 50000},
	{Id: "single_month_unlimited", Title: "1个月(无限)", Amount: 98, Months: 1, CreditPerMonth: 50000},
	{Id: "single_year", Title: "12个月", Amount: 268, OriginAmount: 336, Months: 12, CreditPerMonth: 500},
	{Id: "single_month", Title: "1个月", Amount: 38, Months: 1, CreditPerMonth: 500},
}

func GetPlanById(planId string) *Plan {
	for _, x := range plans {
		if x.Id == planId {
			return &x
		}
	}

	return nil
}

package enums

import (
	"store/pkg/sdk/conv"
	"strings"
	"time"
)

type PaymentCycle string

const (
	PaymentCycleTrialPrefix = "tried"

	PaymentCycle_Annually       PaymentCycle = "annually"
	PaymentCycle_Quarterly      PaymentCycle = "quarterly"
	PaymentCycle_Monthly        PaymentCycle = "monthly"
	PaymentCycle_Daily          PaymentCycle = "daily"
	PaymentCycle_Tried3Monthly  PaymentCycle = PaymentCycleTrialPrefix + "3_monthly"
	PaymentCycle_Tried7Monthly  PaymentCycle = PaymentCycleTrialPrefix + "7_monthly"
	PaymentCycle_Tried14Monthly PaymentCycle = PaymentCycleTrialPrefix + "14_monthly"
)

func (t PaymentCycle) ExpireAt(from time.Time) time.Time {

	day := time.Hour * 24
	month := 31 * day
	year := 365 * day

	switch t {
	case PaymentCycle_Annually:
		//return from.AddDate(1, 0, 0)
		return from.Add(year)
	case PaymentCycle_Quarterly:
		//return from.AddDate(0, 3, 0)
		return from.Add(month * 3)
	case PaymentCycle_Monthly:
		return from.Add(month)
	case PaymentCycle_Daily:
		return from.Add(day)
	case PaymentCycle_Tried3Monthly:
		return from.Add(day * 3)
	case PaymentCycle_Tried7Monthly:
		return from.Add(day * 7)
	case PaymentCycle_Tried14Monthly:
		return from.Add(day * 14)
	}
	return from
}

func (PaymentCycle) Values() []string {
	return []string{
		PaymentCycle_Daily.String(),
		PaymentCycle_Tried3Monthly.String(),
		PaymentCycle_Tried7Monthly.String(),
		PaymentCycle_Tried14Monthly.String(),
		PaymentCycle_Monthly.String(),
		PaymentCycle_Quarterly.String(),
		PaymentCycle_Annually.String(),
	}
}

func (t PaymentCycle) String() string {
	return string(t)
}

func PaymentCycleByCode(code int) PaymentCycle {
	return map[int]PaymentCycle{
		1: PaymentCycle_Daily,
		2: PaymentCycle_Tried3Monthly,
		3: PaymentCycle_Tried7Monthly,
		4: PaymentCycle_Tried14Monthly,
		5: PaymentCycle_Monthly,
		6: PaymentCycle_Quarterly,
		7: PaymentCycle_Annually,
	}[code]
}

func (t PaymentCycle) Name() string {
	return map[PaymentCycle]string{
		PaymentCycle_Daily:          "日付",
		PaymentCycle_Tried14Monthly: "Trial 14",
		PaymentCycle_Tried7Monthly:  "Trial 7",
		PaymentCycle_Tried3Monthly:  "Trial 3",
		PaymentCycle_Monthly:        "Monthly",
		PaymentCycle_Quarterly:      "Quarterly",
		PaymentCycle_Annually:       "Annually",
	}[t]
}

func (t PaymentCycle) Color() string {
	return map[PaymentCycle]string{
		PaymentCycle_Daily:          "#b37feb",
		PaymentCycle_Tried14Monthly: "#b37feb",
		PaymentCycle_Tried7Monthly:  "#b37feb",
		PaymentCycle_Tried3Monthly:  "#b37feb",
		PaymentCycle_Monthly:        "#722ed1",
		PaymentCycle_Quarterly:      "#389e0d",
		PaymentCycle_Annually:       "#eb2f96",
	}[t]
}

func (t PaymentCycle) Sort() int {
	for i, x := range t.Values() {
		if t.String() == x {
			return i
		}
	}
	return 0
}

func (t PaymentCycle) TrialDays() int64 {
	return conv.Int64(strings.Split(t.String(), "_")[0][5:])
}

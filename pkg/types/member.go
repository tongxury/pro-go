package types

import (
	"store/pkg/sdk/helper"
)

type MemberExtra struct {
}

type MemberSubscribeExtra struct {
	OutSubId    string `json:"outSubId"`
	CancelAt    int64  `json:"cancelAt"`
	RenewTimes  int64  `json:"renewTimes"`
	LastRenewAt int64  `json:"lastRenewAt"`
}

func (t MemberSubscribeExtra) Merge(extra MemberSubscribeExtra) MemberSubscribeExtra {

	return MemberSubscribeExtra{
		OutSubId:    helper.OrString(extra.OutSubId, t.OutSubId),
		CancelAt:    helper.OrInt64(extra.CancelAt, t.CancelAt),
		RenewTimes:  helper.OrInt64(extra.RenewTimes, t.RenewTimes),
		LastRenewAt: helper.OrInt64(extra.LastRenewAt, t.LastRenewAt),
	}
}

type Price struct {
	Value    float64 `json:"value,omitempty"`
	Original float64 `json:"original,omitempty"`
	CouponID string  `json:"couponId"`
	Save     string  `json:"save,omitempty"`
}

type Metadata struct {
	Level       string            `json:"level,omitempty"`
	Disable     bool              `json:"disable,omitempty"`
	Prices      map[string]Price  `json:"prices,omitempty"`
	ModelLimits map[string]*Limit `json:"modelLimits,omitempty"`
	//Models         []Model           `json:"models,omitempty"`
	//Functions      []Function        `json:"functions,omitempty"`
	FunctionLimits map[string]*Limit `json:"functionLimits,omitempty"`
	//OtherFunctions []Function        `json:"otherFunctions,omitempty"`
	Label   string `json:"label,omitempty"`
	Suggest bool   `json:"suggest,omitempty"`
}

func (t *Metadata) GetPrice(cycle string) *Price {

	if t == nil {
		return nil
	}

	if p, found := t.Prices[cycle]; found {
		return &p
	}

	return nil
}

func (t Metadata) Copy() *Metadata {
	var c = t
	return &c
}

type Model struct {
	Name    string `json:"name,omitempty"`
	Queries int64  `json:"queries,omitempty"`
	Used    int64  `json:"used,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Label   string `json:"label,omitempty"`
}

type Function struct {
	Name    string `json:"name,omitempty"`
	Queries int64  `json:"queries,omitempty"`
	Used    int64  `json:"used,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

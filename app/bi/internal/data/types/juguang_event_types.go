package types

import "time"

type JuGuangEvent struct {
	ClickId      string    `json:"clickId" ch:"clickId"`
	CaidMd5      string    `json:"caidMd5" ch:"caidMd5"`
	Paid         string    `json:"paid" ch:"paid"`
	Caid         string    `json:"caid" ch:"caid"`
	Ua           string    `json:"ua" ch:"ua"`
	Ua1          string    `json:"ua1" ch:"ua1"`
	IdfaMd5      string    `json:"idfaMd5" ch:"idfaMd5"`
	Os           string    `json:"os" ch:"os"`
	Ip           string    `json:"ip" ch:"ip"`
	Ts           time.Time `json:"ts" ch:"ts"`
	AppId        string    `json:"appId" ch:"appId"`
	AdvertiserId string    `json:"advertiserId" ch:"advertiserId"`
	UnitId       string    `json:"unitId" ch:"unitId"`
	CampaignId   string    `json:"campaignId" ch:"campaignId"`
	CreativityId string    `json:"creativityId" ch:"creativityId"`
	Placement    string    `json:"placement" ch:"placement"`
	AndroidId    string    `json:"androidId" ch:"androidId"`
	Oaid         string    `json:"oaid" ch:"oaid"`
	OaidMd5      string    `json:"oaidMd5" ch:"oaidMd5"`
	Imei         string    `json:"imei" ch:"imei"`
	DeviceId     string    `json:"deviceId" ch:"deviceId"`
	CreatedAt    int64     `json:"createdAt" ch:"createdAt"`
}
type JuGuangEvents []JuGuangEvent

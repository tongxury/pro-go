package conf

import (
	"fmt"
	"golang.org/x/oauth2"
	"store/app/user/internal/data/repo"
	"store/pkg/enums"
	"store/pkg/types"
)

type BizConfig struct {
	Member struct {
		Open              bool
		PromotionExpireAt string
		Models            Descriptions `json:"models,omitempty"`
		Functions         Descriptions `json:"functions,omitempty"`
		MoreFunctions     Descriptions `json:"moreFunctions,omitempty"`
		Metadata          MemberMetadata
		NewUserIdBegin    int64
		NewUserTimeFrom   string
	}
	Banners struct {
		Member       Banner
		Invite       Banner
		Promotion    Banner
		BackToSchool Banner
	}
	Promotions repo.Promotions
	Oauth2     Oauth2
	Payment    Payment
}

type Banner struct {
	ImageUrl string
	Link     string
	Id       string
}

type MemberMetadata struct {
	Free    types.Metadata
	Basic   types.Metadata
	Pro     types.Metadata
	ProPlus types.Metadata
}

func (t *MemberMetadata) Get(level string) (*types.Metadata, error) {
	switch level {
	case enums.MemberLevel_Free.String():
		return t.Free.Copy(), nil
	case enums.MemberLevel_Basic.String():
		return t.Basic.Copy(), nil
	case enums.MemberLevel_Pro.String():
		return t.Pro.Copy(), nil
	case enums.MemberLevel_ProPlus.String():
		return t.ProPlus.Copy(), nil
	}

	return nil, fmt.Errorf("no metadata found by level: %s", level)
}

type Payment struct {
	SuccessUrl string
	CancelUrl  string
}

type Oauth2 struct {
	Redirect     map[string]string
	CookieDomain string `yaml:"cookieDomain"`
	Google       oauth2.Config
}

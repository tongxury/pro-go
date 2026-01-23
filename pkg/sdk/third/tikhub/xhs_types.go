package tikhub

import "strings"

type XHSUser struct {
	Id                        string `json:"id"`
	Avatar                    string `json:"avatar"`
	Fans                      int    `json:"fans"`
	ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
	RedOfficialVerified       bool   `json:"red_official_verified"`
	Self                      bool   `json:"self"`
	Reason                    string `json:"reason"`
	RedOfficialVerifyType     int    `json:"red_official_verify_type"`
	Name                      string `json:"name"`
	Desc                      string `json:"desc"`
	Image                     string `json:"image"`
	Followed                  bool   `json:"followed"`
	RedId                     string `json:"red_id"`
	Link                      string `json:"link"`
	SubTitle                  string `json:"sub_title"`
}

// GetNotes 笔记·200 | 粉丝·2万
func (u XHSUser) GetNotes() string {
	if u.Desc == "" {
		return "0"
	}

	parts := strings.Split(u.Desc, "|")
	if len(parts) != 2 {
		return "0"
	}

	return strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(parts[0]), "笔记·", ""))
}

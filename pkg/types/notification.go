package types

type NotificationExtra struct {
	Title    string    `json:"title"`
	Contents []Content `json:"contents"`
}

type Content struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value"`
}

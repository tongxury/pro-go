package vikingdb

type Response[T any] struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Result    T      `json:"result"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	ImageTokens      int `json:"image_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

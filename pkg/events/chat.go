package events

type ChatEvent struct {
	Timestamp    int64
	UserID       int64  `json:"user_id"`
	FunctionName string `json:"function_name"`
	Model        string
	Url          string
	Query        string
	Image        string
	Answer       string
	Status       string
}

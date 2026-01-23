package eventsource

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
)

type Chunk struct {
	Id      string `json:"id,omitempty"`
	Data    any    `json:"data,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewChunk(id string, code int, data any, message string) *Chunk {
	return &Chunk{
		Id:      id,
		Data:    data,
		Code:    code,
		Message: message,
	}
}

func (t *Chunk) String() string {
	s, err := json.Marshal(t)

	if err != nil {
		log.Errorw("Marshal err", err, "t", t)
		return ""
	}

	if len(s) == 0 || string(s) == "" {
		return "data: \n\n"
	}

	d := "data: " + string(s) + "\n\n"

	return d
}

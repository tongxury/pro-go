package bitquery

import "time"

type Block struct {
	Date       string    `json:"Date,omitempty"`
	Hash       string    `json:"Hash,omitempty"`
	Height     string    `json:"Height,omitempty"`
	ParentHash string    `json:"ParentHash,omitempty"`
	ParentSlot string    `json:"ParentSlot,omitempty"`
	Slot       string    `json:"Slot,omitempty"`
	Time       time.Time `json:"Time,omitempty"`
}

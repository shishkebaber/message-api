package model

type Message struct {
	Name     string `json:"name"  validate:"required"`
	Payload  string `json:"payload" validate:"required"`
	Metadata `json:"metadata" validate:"required"`
}

type Metadata struct {
	Timestamp string `json:"timestamp" validate:"required,timestamp"`
	Uuid      string `json:"uuid" validate:"required,uuid"`
}

type MessageKey struct{}

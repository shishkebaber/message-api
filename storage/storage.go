package storage

import "github.com/shishkebaber/message-api/model"

type Storage interface {
	CreateMessage(message *model.Message) error
	FetchAllMessages() ([]*model.Message, error)
	FetchMessage(name string) (*model.Message, error)
}

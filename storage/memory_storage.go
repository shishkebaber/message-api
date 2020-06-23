package storage

import (
	"fmt"
	"github.com/shishkebaber/message-api/model"
	"github.com/sirupsen/logrus"
	"sync"
)

var ErrMessageNotFound = fmt.Errorf("Message not found")
var ErrMessageWithNameExist = fmt.Errorf("Message with such name already exist")
var ErrMessageWithUUIDExist = fmt.Errorf("Message with such UUID already exist")

type MStorage struct {
	mutex    sync.Mutex
	Messages map[string]*model.Message
	Logger   *logrus.Logger
}

func NewMStorage() *MStorage {
	m := sync.Mutex{}
	l := logrus.New()
	return &MStorage{mutex: m, Messages: make(map[string]*model.Message), Logger: l}
}

func (s *MStorage) checkUUIDUnique(uuid string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range s.Messages {
		if uuid == v.Uuid {
			return true
		}
	}
	return false
}

func (s *MStorage) checkUniqueFields(message *model.Message) error {
	/*s.mutex.Lock()
	defer s.mutex.Unlock()*/
	if _, ok := s.Messages[message.Name]; ok {
		return ErrMessageWithNameExist
	} else if s.checkUUIDUnique(message.Uuid) {
		return ErrMessageWithUUIDExist
	}
	return nil
}

func (s *MStorage) CreateMessage(message *model.Message) error {
	s.Logger.Info("Inserting message")
	err := s.checkUniqueFields(message)
	if err != nil {
		s.Logger.Error("Error during insert: ", err)
		return err
	}
	s.mutex.Lock()
	s.Messages[message.Name] = message
	s.mutex.Unlock()
	return nil
}

func (s *MStorage) FetchAllMessages() ([]*model.Message, error) {
	msgs := []*model.Message{}
	s.mutex.Lock()
	s.Logger.Info("Fetching all Messages from memory")
	for _, v := range s.Messages {
		msgs = append(msgs, v)
	}
	s.mutex.Unlock()
	return msgs, nil
}

func (s *MStorage) FetchMessage(name string) (*model.Message, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Logger.Info("Fetching single message from memory")
	if msg, ok := s.Messages[name]; !ok {
		return nil, ErrMessageNotFound
	} else {
		return msg, nil
	}
}

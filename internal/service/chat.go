package service

import (
	"github.com/bopoh24/realty-bot/internal/models"
)

type ChatStore interface {
	Load(map[models.ChatID]models.Chat) error
	Save(map[models.ChatID]models.Chat) error
}

type ChatService struct {
	chats map[models.ChatID]models.Chat
	store ChatStore
}

// NewChatService returns new subscribers instance.
func NewChatService(store ChatStore) (*ChatService, error) {
	s := &ChatService{
		chats: make(map[models.ChatID]models.Chat),
		store: store,
	}
	if err := s.store.Load(s.chats); err != nil {
		return nil, err
	}
	return s, nil
}

// Exists checks if subscriber exists.
func (s *ChatService) Exists(chatID models.ChatID) bool {
	_, ok := s.chats[chatID]
	return ok
}

func (s *ChatService) Save(chat models.Chat) error {
	s.chats[chat.ChatID] = chat
	return s.store.Save(s.chats)
}

func (s *ChatService) Delete(chatID models.ChatID) error {
	delete(s.chats, chatID)
	return s.store.Save(s.chats)
}

// List returns list of subscribers.
func (s *ChatService) List() []models.Chat {
	list := make([]models.Chat, 0, len(s.chats))
	for _, chat := range s.chats {
		list = append(list, chat)
	}
	return list
}

package service

import (
	"github.com/bopoh24/realty-bot/internal/models"
)

type UserStore interface {
	Load(map[int64]models.User) error
	Save(map[int64]models.User) error
}

type UserService struct {
	users map[int64]models.User
	store UserStore
}

// NewUserService returns new subscribers instance
func NewUserService(store UserStore) (*UserService, error) {
	s := &UserService{
		users: make(map[int64]models.User),
		store: store,
	}
	if err := s.store.Load(s.users); err != nil {
		return nil, err
	}
	return s, nil
}

// Exists checks if subscriber exists
func (s *UserService) Exists(chatID int64) bool {
	_, ok := s.users[chatID]
	return ok
}

func (s *UserService) Save(user models.User) error {
	s.users[user.ChatID] = user
	return s.store.Save(s.users)
}

func (s *UserService) Delete(chatID int64) error {
	delete(s.users, chatID)
	return s.store.Save(s.users)
}

// List returns list of subscribers
func (s *UserService) List() []models.User {
	list := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		list = append(list, user)
	}
	return list
}

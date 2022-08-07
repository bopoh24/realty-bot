package filestore

import (
	"encoding/json"
	"errors"
	"github.com/bopoh24/realty-bot/internal/models"
	"os"
)

type ChatStore struct {
	filename string
}

// NewChatStore return store instance
func NewChatStore(filename string) *ChatStore {
	return &ChatStore{filename: filename}
}

func (u *ChatStore) Load(output map[models.ChatID]models.Chat) error {
	if _, err := os.Stat(u.filename); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	data, err := os.ReadFile(u.filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &output); err != nil {
		return err
	}
	return nil
}

func (u *ChatStore) Save(users map[models.ChatID]models.Chat) error {
	data, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(u.filename, data, 0644); err != nil {
		return err
	}
	return nil
}

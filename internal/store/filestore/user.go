package filestore

import (
	"encoding/json"
	"errors"
	"os"
	"realty_bot/internal/models"
)

type UserStore struct {
	filename string
}

// NewUserStore return store instance
func NewUserStore(filename string) *UserStore {
	return &UserStore{filename: filename}
}

func (u *UserStore) Load(output map[int64]models.User) error {
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

func (u *UserStore) Save(users map[int64]models.User) error {
	data, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(u.filename, data, 0644); err != nil {
		return err
	}
	return nil
}

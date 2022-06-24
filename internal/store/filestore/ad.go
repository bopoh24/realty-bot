package filestore

import (
	"encoding/json"
	"errors"
	"github.com/bopoh24/realty-bot/internal/models"
	"os"
)

type AdStore struct {
	filename string
}

// NewAdStore returns file store for ads
func NewAdStore(filename string) *AdStore {
	return &AdStore{
		filename: filename,
	}
}

func (a *AdStore) Save(ads []models.Ad) error {
	data, err := json.MarshalIndent(ads, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(a.filename, data, 0644); err != nil {
		return err
	}
	return nil
}

func (a *AdStore) Load() (map[string]models.Ad, error) {
	data, err := os.ReadFile(a.filename)
	if err != nil {
		return nil, err
	}
	var sentAds []models.Ad
	if err = json.Unmarshal(data, &sentAds); err != nil {
		return nil, err
	}
	sentAdsMap := make(map[string]models.Ad)
	for _, ad := range sentAds {
		sentAdsMap[ad.Link] = ad
	}
	return sentAdsMap, nil
}

func (a *AdStore) IsEmpty() bool {
	_, err := os.Stat(a.filename)
	return errors.Is(err, os.ErrNotExist)
}

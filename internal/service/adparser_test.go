package service

import (
	"github.com/stretchr/testify/assert"
	"realty_bot/internal/store/filestore"
	"realty_bot/pkg/log"
	"testing"
)

const searchURL = "https://www.bazaraki.com/real-estate/houses-and-villas-rent/?city_districts=5732&city_districts=5815&city_districts=5772&city_districts=5771&city_districts=5776&city_districts=5735&city_districts=5774&city_districts=5770&city_districts=5773&city_districts=5733&city_districts=5734&city_districts=5736&city_districts=5545&city_districts=5503&price_max=1500"

func TestLoadBodyOK(t *testing.T) {
	b := NewAdParseService("https://google.com", log.NewLogger(), nil)
	body, err := b.loadHTML()
	assert.NoError(t, err)
	body.Close()
}

func TestLoadBodyErrRequest(t *testing.T) {
	b := NewAdParseService("https://fdfadfsdfasdh.com", log.NewLogger(), nil)
	body, err := b.loadHTML()
	assert.Error(t, err)
	assert.Nil(t, body)
}

func TestAds(t *testing.T) {
	b := NewAdParseService(searchURL, log.NewLogger(), nil)
	body, err := b.loadHTML()
	assert.NoError(t, err)
	ads, err := b.parseAds(body)
	assert.NoError(t, err)
	assert.NotEmpty(t, ads)
	t.Logf("%#v", ads[0])
}

func TestNewAds(t *testing.T) {
	b := NewAdParseService(searchURL, log.NewLogger(), filestore.NewAdStore("ads_test.json"))
	ads, err := b.NewAds()
	assert.NoError(t, err)
	assert.NotEmpty(t, ads)
	t.Logf("%#v", ads[0])
}

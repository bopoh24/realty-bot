package service

import (
	"github.com/bopoh24/realty-bot/internal/store/filestore"
	"github.com/bopoh24/realty-bot/pkg/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const htmlSampleFile = "../../testing/html_sample.html"
const storeFile = "ads_test.json"

func TestBodyDataOK(t *testing.T) {
	b, err := NewAdParseService("https://google.com", log.NewLogger(), nil)
	assert.NoError(t, err)

	body, err := b.bodyData()
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}

func TestBodyDataErrRequest(t *testing.T) {
	b, err := NewAdParseService("https://non-existent.domain", log.NewLogger(), nil)
	assert.NoError(t, err)

	body, err := b.bodyData()
	assert.Error(t, err)
	assert.Nil(t, body)
}

func TestAds(t *testing.T) {
	b, err := NewAdParseService("https://non-existent.domain", log.NewLogger(), nil)
	assert.NoError(t, err)
	body, err := ioutil.ReadFile(htmlSampleFile)
	assert.NoError(t, err)
	ads, err := b.parsedAds(body)
	assert.NoError(t, err)
	assert.NotEmpty(t, ads)
	ad := ads[0]
	assert.NotEmpty(t, ad.Link)
	assert.NotEmpty(t, ad.Title)
	assert.NotEmpty(t, ad.Price)
	assert.NotEmpty(t, ad.Location)
	assert.False(t, ad.Datetime.IsZero())
	// only 4 ads before "Ads from other regions"
	assert.Len(t, ads, 4)
}

func TestNewAds(t *testing.T) {

	b, err := NewAdParseService("https://non-existent.domain",
		log.NewLogger(), filestore.NewAdStore(storeFile))
	assert.NoError(t, err)
	assert.True(t, b.store.IsEmpty())
	defer func() {
		_ = os.Remove(storeFile)
	}()

	body, err := ioutil.ReadFile(htmlSampleFile)
	assert.NoError(t, err)

	ads, err := b.parsedAds(body)
	assert.NoError(t, err)
	assert.NotEmpty(t, ads)

	newAds, err := b.newAds(ads[1:])
	assert.NoError(t, err)
	assert.Empty(t, newAds)

	newAds, err = b.newAds(ads)
	assert.NoError(t, err)
	assert.Len(t, newAds, 1)

}

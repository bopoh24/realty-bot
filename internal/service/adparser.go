package service

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bopoh24/realty-bot/internal/models"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type AdStoreInterface interface {
	Save(ads []models.Ad) error
	Load() (map[models.AdLink]models.Ad, error)
	IsEmpty() bool
}

type AdParseService struct {
	parseURL *url.URL
	store    AdStoreInterface
	logger   *zerolog.Logger
}

// NewAdParseService returns instance of service.
func NewAdParseService(searchLink string, logger *zerolog.Logger, store AdStoreInterface) (*AdParseService, error) {
	parseURL, err := url.Parse(searchLink)
	if err != nil {
		return nil, err
	}
	return &AdParseService{
		parseURL: parseURL,
		logger:   logger,
		store:    store,
	}, nil
}

func (b *AdParseService) bodyData() ([]byte, error) {
	res, err := http.Get(b.parseURL.String())
	if err != nil {
		return nil, fmt.Errorf("can't get search link: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			b.logger.Error().Msg(err.Error())
		}
	}()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("can't load html body, status=%d", res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

func (b *AdParseService) parsedAds(bodyData []byte) ([]models.Ad, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyData))
	if err != nil {
		return nil, err
	}
	result := make([]models.Ad, 0)

	otherAdsBlockFound := false
	otherAdsBlockText := "Ads from other regions"

	doc.Find(".list-simple__output").Children().Each(func(i int, s *goquery.Selection) {
		if s.Is("h2.header") && strings.TrimSpace(s.Text()) == otherAdsBlockText {
			otherAdsBlockFound = true
			return
		}
		if otherAdsBlockFound {
			return
		}
		if s.Is("li.announcement-container") {
			b.parseAdBlock(s, &result)
		}
	})
	return result, nil
}

func (b *AdParseService) parseAdBlock(blockSelection *goquery.Selection, ads *[]models.Ad) {
	var err error
	blockSelection.Find(".list-announcement-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		titleLink := s.Find(".announcement-block__title")
		title := strings.TrimSpace(titleLink.Text())
		link, exist := titleLink.Attr("href")
		if !exist {
			b.logger.Error().Msgf("unable to get link for %s", title)
			return
		}
		link = b.parseURL.Scheme + "://" + path.Join(b.parseURL.Host, strings.TrimSpace(link))

		blockDateText := strings.TrimSpace(s.Find(".announcement-block__date").Text())
		blockDateSplit := strings.Split(blockDateText, ",")
		datetime := b.parseDate(blockDateSplit[0])
		location := strings.TrimSpace(blockDateSplit[len(blockDateSplit)-1])
		price := 0
		s.Find("meta").Each(func(i int, m *goquery.Selection) {
			val, exists := m.Attr("itemprop")
			if !exists {
				return
			}
			if val == "price" {
				priceStr, exist := m.Attr("content")
				if !exist {
					b.logger.Error().Msgf("unable to get price for %s", title)
				}
				price, err = strconv.Atoi(strings.Split(priceStr, ".")[0])
				if err != nil {
					b.logger.Error().Msgf("unable to convert price for %s", title)
				}
			}
		})
		ad := models.Ad{
			Title:    title,
			Link:     models.AdLink(link),
			Location: location,
			Price:    price,
			Datetime: datetime,
		}
		*ads = append(*ads, ad)
	})
}

func (b *AdParseService) parseDate(dateStr string) time.Time {
	dateStrSplit := strings.Split(dateStr, " ")
	now := time.Now()
	switch dateStrSplit[0] {
	case "Yesterday":
		t, err := time.Parse("15:04", dateStrSplit[1])
		if err != nil {
			b.logger.Error().Msgf("unable to parse time for Yesterday: %s", err)
			return time.Time{}
		}
		return time.Date(now.Year(), now.Month(), now.Day()-1, t.Hour(), t.Minute(), 0, 0, now.Location())
	case "Today":
		t, err := time.Parse("15:04", dateStrSplit[1])
		if err != nil {
			b.logger.Error().Msgf("unable to parse time for Today: %s", err)
			return time.Time{}
		}
		return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	default:
		t, err := time.Parse("02.01.2006 15:04", dateStr)
		if err != nil {
			b.logger.Error().Msgf("unable to parse time: %s", err)
			return time.Time{}
		}
		return t
	}
}

func (b *AdParseService) newAds(parsedAds []models.Ad) ([]models.Ad, error) {
	if b.store.IsEmpty() {
		if err := b.store.Save(parsedAds); err != nil {
			return nil, err
		}
		// create initial ads
		return nil, nil
	}
	// read store and compare
	sentAdsMap, err := b.store.Load()
	// compare
	newAds := make([]models.Ad, 0)
	for _, newAd := range parsedAds {
		if _, ok := sentAdsMap[newAd.Link]; ok {
			continue
		}
		newAds = append(newAds, newAd)
	}
	if len(newAds) != 0 {
		if err = b.store.Save(parsedAds); err != nil {
			return nil, err
		}
	}
	return newAds, nil
}

// AdsToNotify returns new ads
func (b *AdParseService) AdsToNotify() ([]models.Ad, error) {
	bodyData, err := b.bodyData()
	if err != nil {
		return nil, err
	}
	ads, err := b.parsedAds(bodyData)
	if err != nil {
		return nil, err
	}
	return b.newAds(ads)
}

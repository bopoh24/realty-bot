package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"realty_bot/internal/models"
	"strconv"
	"strings"
	"time"
)

type AdStoreInterface interface {
	Save(ads []models.Ad) error
	Load() (map[string]models.Ad, error)
	IsEmpty() bool
}

type AdParseService struct {
	searchLink string
	store      AdStoreInterface
	logger     *zerolog.Logger
}

// NewAdParseService returns instance of service
func NewAdParseService(searchLink string, logger *zerolog.Logger, store AdStoreInterface) *AdParseService {
	return &AdParseService{
		searchLink: searchLink,
		logger:     logger,
		store:      store,
	}
}

func (b *AdParseService) loadHTML() (io.ReadCloser, error) {
	res, err := http.Get(b.searchLink)
	if err != nil {
		return nil, fmt.Errorf("can't get search link: %w", err)
	}

	if res.StatusCode != 200 {
		if err = res.Body.Close(); err != nil {
			b.logger.Error().Msg(err.Error())
		}
		return nil, fmt.Errorf("can't load html body, status=%d", res.StatusCode)
	}
	return res.Body, nil
}

func (b *AdParseService) parseAds(body io.ReadCloser) ([]models.Ad, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	result := make([]models.Ad, 0)
	doc.Find(".list-announcement-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		titleLink := s.Find(".announcement-block__title")
		title := strings.TrimSpace(titleLink.Text())
		link, exist := titleLink.Attr("href")
		if !exist {
			b.logger.Error().Msgf("unable to get link for %s", title)
			return
		}
		link = "https://www.bazaraki.com" + strings.TrimSpace(link)

		blockDateText := strings.TrimSpace(s.Find(".announcement-block__date").Text())
		blockDateSplit := strings.Split(blockDateText, ",")
		datetime := b.parseDate(blockDateSplit[0])
		location := strings.TrimSpace(blockDateSplit[len(blockDateSplit)-1])
		if !strings.HasPrefix(location, "Larnaka") {
			return
		}
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
			Link:     link,
			Location: location,
			Price:    price,
			Datetime: datetime,
		}
		result = append(result, ad)
	})
	return result, nil
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

// NewAds returns new ads
func (b *AdParseService) NewAds() ([]models.Ad, error) {
	htmlReader, err := b.loadHTML()
	if err != nil {
		return nil, err
	}
	defer htmlReader.Close()

	ads, err := b.parseAds(htmlReader)
	if err != nil {
		return nil, err
	}
	if b.store.IsEmpty() {
		if err = b.store.Save(ads); err != nil {
			return nil, err
		}
		// create initial ads
		return nil, nil
	}
	// read store and compare
	sentAdsMap, err := b.store.Load()
	// compare
	adsToSend := make([]models.Ad, 0)
	for _, newAd := range ads {
		if _, ok := sentAdsMap[newAd.Link]; ok {
			continue
		}
		adsToSend = append(adsToSend, newAd)
	}
	if len(adsToSend) != 0 {
		if err = b.store.Save(ads); err != nil {
			return nil, err
		}
	}
	return adsToSend, nil
}

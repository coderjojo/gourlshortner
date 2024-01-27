package shortner

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type urlEntry struct {
	orignalURL     string
	expirationTime time.Time
	creationTime   time.Time
}

type UrlShortner struct {
	urlMapping map[string]urlEntry
	expiration time.Duration
	Logger     *log.Logger
	mutex      sync.Mutex
}

func NewUrlShorter(expiration time.Duration, logger *log.Logger) *UrlShortner {

	return &UrlShortner{
		urlMapping: make(map[string]urlEntry),
		expiration: expiration,
		Logger:     logger,
		mutex:      sync.Mutex{},
	}
}

func (us *UrlShortner) Redirect(shortUrl string) (string, error) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	entry, ok := us.urlMapping[shortUrl]
	if !ok || entry.expirationTime.Before(time.Now()) {
		us.Logger.Printf("Redirect failed: %s not found or expired", shortUrl)
		return "", errors.New("Short URL not found or expired")
	}

	us.Logger.Printf("Redirecting : %s -> %s\n", shortUrl, entry.orignalURL)
	return entry.orignalURL, nil
}

func (us *UrlShortner) ShortenURL(orignalURL string) (string, error) {

	us.mutex.Lock()
	defer us.mutex.Unlock()

	// check if URL already exists
	for key, entry := range us.urlMapping {
		if entry.orignalURL == orignalURL && entry.expirationTime.After(time.Now()) {
			us.Logger.Printf("URL already shortened: %s -> %s\n", orignalURL, key)
			return key, nil
		}
	}

	shortURL := generateURL()

	us.urlMapping[shortURL] = urlEntry{
		orignalURL:     orignalURL,
		expirationTime: time.Now().Add(us.expiration),
		creationTime:   time.Now(),
	}

	us.Logger.Printf("URL shortened: %s -> %s\n", orignalURL, shortURL)
	return shortURL, nil
}

// NOTE: Improve the logic
func generateURL() string {
	return "short_" + fmt.Sprint(time.Now().UnixNano())[:8]
}

func (us *UrlShortner) URLStats() string {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	stats := make(map[string]interface{})

	for key, value := range us.urlMapping {
		stats[key] = map[string]interface{}{
			"orignalURL": value.orignalURL,
			"expiration": value.expirationTime.Format(time.RFC3339),
			"created":    value.creationTime.Format(time.RFC3339),
		}
	}

	statsJson, _ := json.MarshalIndent(stats, "", " ")

	return string(statsJson)
}

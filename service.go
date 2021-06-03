package main

import (
	"log"
	"strconv"

	gcache "github.com/patrickmn/go-cache"
)

type BusinessLogic struct {
	db Repository
	c  gcache.Cache
}

func (d BusinessLogic) GenerateShortURLCode(longURL string) (string, error) {

	id, err := d.db.Save(longURL)
	if err != nil {
		return "", err
	}

	newURL := longURL + strconv.FormatInt(id, 10)

	shortURLCode := substring(computeSHA256Base64(newURL))

	go d.updateShortURLCode(id, shortURLCode) // update in background

	return shortURLCode, nil
}

func (d BusinessLogic) updateShortURLCode(id int64, shortURL string) {

	_, err := d.db.Update(id, shortURL)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

}

func (d BusinessLogic) search(shortURL string) (string, error) {

	// Check first for key in cache
	if x, found := d.c.Get(shortURL); found {
		cacheShortUrl := *x.(**URLInfo)
		log.Printf("INFO: key: %s was read from cache Value: %s", shortURL, cacheShortUrl.OriginalURL)
		return cacheShortUrl.OriginalURL, nil
	}

	// Key not found in cache, search database
	urlInfo, err := d.db.SearchByShortURL(shortURL)
	if err != nil {
		return "", err
	}

	// Store new item in cache
	d.c.Set(shortURL, &urlInfo, gcache.DefaultExpiration)

	return urlInfo.OriginalURL, nil

}

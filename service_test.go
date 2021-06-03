package main

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

type mockRepository struct {
	id int64
}

func (r *mockRepository) Save(longURL string) (int64, error) {

	return r.id + 1, nil
}

func (r *mockRepository) Update(id int64, shortURL string) (int64, error) {
	return 0, nil
}

func (r *mockRepository) SearchByShortURL(shortURL string) (*URLInfo, error) {

	return &URLInfo{
		OriginalURL: "https://www.golangprograms.com/data-structure-and-algorithms.html",
		ShortenURL:  "YjEwMj",
		ID:          1,
	}, nil
}

// instantiate datastore
var testDataStore mockRepository

// instatiate cache
// Create a cache with a default expiration time of 5 minutes, and which
// purges expired items every 10 minutes
var testCache = cache.New(5*time.Minute, 10*time.Minute)
var serviceTest = BusinessLogic{&testDataStore, *testCache}

func TestGenerateShortURLCode(t *testing.T) {

	url := `https://www.golangprograms.com/data-structure-and-algorithms.html`

	expected := "YjEwMj"

	shortURLCode, _ := serviceTest.GenerateShortURLCode(url)

	if shortURLCode != expected {
		t.Errorf("want: %s but got: %s", expected, shortURLCode)
	}

}

func TestSearchByShortURL(t *testing.T) {

	expected := `https://www.golangprograms.com/data-structure-and-algorithms.html`

	shortURLCode := "YjEwMj"

	urlInfo, _ := serviceTest.db.SearchByShortURL(shortURLCode)

	if urlInfo.OriginalURL != expected {
		t.Errorf("want: %s but got: %s", expected, urlInfo.OriginalURL)
	}

}

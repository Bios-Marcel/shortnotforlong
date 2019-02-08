package tests

import (
	"net/http"
	"testing"

	linkshortener "github.com/Bios-Marcel/shortnotforlong"
)

func TestShortener(t *testing.T) {
	shortener := linkshortener.NewShortener(55555)
	go func() {
		shortenerError := shortener.Start()
		if shortenerError != nil {
			t.Error("Error starting shortener:", shortenerError)
		}
	}()

	{
		URL := "https://www.google.com/"
		shortenedURL := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}
		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}
	}

	{
		URL := "https://duckduckgo.com/"
		shortenedURL := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}
		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}
	}
}

package tests

import (
	linkshortener "github.com/Bios-Marcel/shortnotforlong"
	"net/http"
	"strings"
	"testing"
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

	{
		suffix := ".png"
		URL := "https://duckduckgo.com/img" + suffix
		shortenedURL := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}

		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}

		if !strings.HasSuffix(shortenedURL, suffix) {
			t.Errorf("Shortened url does not end with: %s", suffix)
		}
	}

}

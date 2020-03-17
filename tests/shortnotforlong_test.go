package tests

import (
	linkshortener "github.com/Bios-Marcel/shortnotforlong"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// URLs without a colon prefix (e.g. http://) are not supported by the go http library

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

	{
		suffix := ".webm"
		URL := "http://duckduckgo.com/vid" + suffix
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

	{
		URL := "https://github.com"
		shortenedURL := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}

		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}

		var urlSuffix = regexp.MustCompile(`.*(/\d+)$`)
		var matches = urlSuffix.FindStringSubmatch(shortenedURL)
		if len(matches) == 0 {
			t.Errorf("Shortened url does not end with '/{number}' - %s", shortenedURL)
		}
	}

}

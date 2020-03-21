package tests

import (
	"net/http"
	"regexp"
	"testing"

	linkshortener "github.com/Bios-Marcel/shortnotforlong"
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
		shortenedURL, suffix := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}
		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}
		if suffix != "" {
			t.Errorf("A suffix was found where none was present: %s in %s", suffix, URL)
		}
	}

	{
		URL := "https://duckduckgo.com/"
		shortenedURL, suffix := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}
		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}
		if suffix != "" {
			t.Errorf("A suffix was found where none was present: %s in %s", suffix, URL)
		}
	}

	{
		URL := "https://duckduckgo.com/img.png"
		shortenedURL, suffix := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}

		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}

		if suffix != ".png" {
			t.Errorf("Incorrect suffix %s found in %s", suffix, URL)
		}
	}

	{
		URL := "http://duckduckgo.com/vid.webm"
		shortenedURL, suffix := shortener.Shorten(URL)
		response, httpError := http.Head(shortenedURL)
		if httpError != nil {
			t.Fatalf("Error shortening looking up: '%s' (%s). Error: %s", shortenedURL, URL, httpError.Error())
		}

		if response.Request.URL.String() != URL {
			t.Errorf("URL was incorrect: %s", response.Request.URL.String())
		}

		if suffix != ".webm" {
			t.Errorf("Incorrect suffix %s found in %s", suffix, URL)
		}
	}

	{
		URL := "https://github.com"
		shortenedURL, suffix := shortener.Shorten(URL)
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

		if suffix != "" {
			t.Errorf("Unexpected suffix %s found in %s", suffix, URL)
		}
	}

}

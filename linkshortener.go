package linkshortener

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// Shortener offers a function to shorten a URL and redirect to the shortened
// URL as soon as a request comes in.
type Shortener struct {
	nextFreeIndex uint16
	shortenedUrls map[uint16]string
	port          int
	httpServer    *http.Server
}

// Shorten takes a url and returns a shortend version that redirects via the
// local webserver.
func (shortener *Shortener) Shorten(url string) string {
	var suffix = ""
	var urlSuffix = regexp.MustCompile(`(http(s)?:\/\/)?.*(\.\w+)+\/.*(\.\w{1,4})$`)
	var matches = urlSuffix.FindStringSubmatch(url)
	if len(matches) > 0 {
		suffix = matches[len(matches) - 1]
	}

	for id, address := range shortener.shortenedUrls {
		if address == url {
			return fmt.Sprintf("http://localhost:%d/%d"+suffix, shortener.port, id)
		}
	}

	newID := shortener.generateID()
	shortener.shortenedUrls[newID] = url

	return fmt.Sprintf("http://localhost:%d/%d"+suffix, shortener.port, newID)
}

func (shortener *Shortener) generateID() uint16 {
	if shortener.nextFreeIndex >= math.MaxUint16 {
		shortener.nextFreeIndex = 0
	} else {
		shortener.nextFreeIndex++
	}

	return shortener.nextFreeIndex
}

//RedirectHandler handles all the redirects for the Server.
type RedirectHandler struct {
	Shortener *Shortener
}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var id = ""
	var idRegex = regexp.MustCompile(`/(\d*)`)
	var matches = idRegex.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 {
		id = matches[1]
	}

	idAsInt, convertError := strconv.ParseUint(id, 10, 16)
	if convertError != nil {
		http.NotFound(w, r)
	} else {
		url, contains := h.Shortener.shortenedUrls[uint16(idAsInt)]
		if contains {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			http.NotFound(w, r)
		}
	}
}

// Close closes the internal http server.
func (shortener *Shortener) Close() {
	shortener.httpServer.Shutdown(context.Background())
}

//NewShortener creates a new server that uses the given port.
func NewShortener(port int) *Shortener {
	shortener := &Shortener{
		shortenedUrls: make(map[uint16]string),
		port:          port,
	}

	handler := RedirectHandler{
		Shortener: shortener,
	}

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	shortener.httpServer = httpServer

	return shortener
}

//Start servers the internal http server, blocks and returns an error on
//failure.
func (shortener *Shortener) Start() error {
	return shortener.httpServer.ListenAndServe()
}

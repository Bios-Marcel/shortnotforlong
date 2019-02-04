package linkshortener

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Shortener offers a function to shorten a URL and redirect to the shortened
// URL as soon as a request comes in.
type Shortener struct {
	shortenedUrls map[uint16]string
	port          int
	httpServer    *http.Server
}

// Shorten takes a url and returns a shortend version that redirects via the
// local webserver.
func (shortener *Shortener) Shorten(url string) string {
	for id, address := range shortener.shortenedUrls {
		if address == url {
			return fmt.Sprintf("http://localhost:%d/%d", shortener.port, id)
		}
	}

	newID := shortener.generateID()
	shortener.shortenedUrls[newID] = url

	return fmt.Sprintf("http://localhost:%d/%d", shortener.port, newID)
}

func (shortener *Shortener) generateID() uint16 {
	rand.Seed(time.Now().UTC().UnixNano())
	randomInt := uint16(rand.Int31n(math.MaxUint16))
	_, contains := shortener.shortenedUrls[randomInt]
	if contains {
		return shortener.generateID()
	}

	return randomInt
}

//RedirectHandler handles all the redirects for the Server.
type RedirectHandler struct {
	Shortener *Shortener
}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	idAsInt, convertError := strconv.ParseUint(id, 10, 16)
	if convertError != nil {
		http.NotFound(w, r)
	} else {
		url, contains := h.Shortener.shortenedUrls[uint16(idAsInt)]
		if contains {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			http.NotFound(w, r)
		}
	}
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

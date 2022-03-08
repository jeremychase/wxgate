package requestlog

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Requestlog wraps a handler to create middleware.
type Requestlog struct {
	handler  http.Handler
	filename string
}

// ServeHTTP writes a request log to file.
func (l *Requestlog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	url := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.Query().Encode())

	l.handler.ServeHTTP(w, r)

	f, err := os.OpenFile(l.filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("error: unable to open log file: %s", err)
		return
	}
	defer func() {
		err := f.Sync()
		if err != nil {
			log.Printf("error: unable to sync log file: %s", err)
		}
		err = f.Close()
		if err != nil {
			log.Printf("error: unable to close log file: %s", err)
		}
	}()

	reqlog := log.New(f, "", log.LstdFlags)
	reqlog.Printf("%s %s %v", r.Method, url, time.Since(start))
}

// NewReqLog returns a Requestlog middleware.
func NewReqLog(mux http.Handler, filename string) *Requestlog {
	return &Requestlog{handler: mux, filename: filename}
}

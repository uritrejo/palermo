package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Received a request: \"%s\" from %s on %s", r.URL.String(), r.RemoteAddr, time.Now().Format(time.RFC822Z))

		// todo: add date & time, r.RemoteAddr
		next.ServeHTTP(w, r)
	})
}

// todo: maybe do a security one
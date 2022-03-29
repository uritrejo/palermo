package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusContinue)
	})

	handler := LoggingMiddleware(fn)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/hello", nil))

	assert.Equal(t, http.StatusContinue, rr.Code)
}

func TestRecoveryMiddleware(t *testing.T) {
	panicFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("tragedy")
	})
	handler := RecoveryMiddleware(panicFn)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, nil)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

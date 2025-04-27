package http

import (
	"net/http"

	"github.com/KrllF/Cloud/consts"
)

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(consts.Attempts).(int); ok {
		return attempts
	}

	return 1
}

// GetRetryFromContext returns the retries for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(consts.Retry).(int); ok {
		return retry
	}

	return 0
}

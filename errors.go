package noxaeapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

// APIError is returned for any non-2xx response from a NoxAeApi server.
// Use errors.As to check for it, or the Is* helper functions below.
type APIError struct {
	Status       int
	Method       string
	Path         string
	Body         json.RawMessage
	RetryAfterMs *int64 // populated for 429 responses when the server sends a Retry-After header
}

func (e *APIError) Error() string {
	switch e.Status {
	case 401:
		return fmt.Sprintf("noxaeapi: unauthorized (%s %s): the API key was missing or invalid", e.Method, e.Path)
	case 403:
		return fmt.Sprintf("noxaeapi: forbidden (%s %s): this API key does not have permission to call this endpoint", e.Method, e.Path)
	case 404:
		return fmt.Sprintf("noxaeapi: not found (%s %s)", e.Method, e.Path)
	case 429:
		if e.RetryAfterMs != nil {
			return fmt.Sprintf("noxaeapi: rate limited (%s %s), retry after %dms", e.Method, e.Path, *e.RetryAfterMs)
		}
		return fmt.Sprintf("noxaeapi: rate limited (%s %s)", e.Method, e.Path)
	default:
		if e.Status >= 500 {
			return fmt.Sprintf("noxaeapi: server error %d (%s %s)", e.Status, e.Method, e.Path)
		}
		return fmt.Sprintf("noxaeapi: unexpected status %d (%s %s)", e.Status, e.Method, e.Path)
	}
}

// NetworkError is returned when a request never completed at all
// (DNS failure, connection refused, timeout). Unlike APIError, there was
// no HTTP response to inspect.
type NetworkError struct {
	Method string
	Path   string
	Err    error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("noxaeapi: network error on %s %s: %v", e.Method, e.Path, e.Err)
}

func (e *NetworkError) Unwrap() error { return e.Err }

// IsUnauthorized reports whether err is an APIError with status 401.
func IsUnauthorized(err error) bool { return hasStatus(err, 401) }

// IsForbidden reports whether err is an APIError with status 403.
func IsForbidden(err error) bool { return hasStatus(err, 403) }

// IsNotFound reports whether err is an APIError with status 404.
func IsNotFound(err error) bool { return hasStatus(err, 404) }

// IsRateLimited reports whether err is an APIError with status 429.
func IsRateLimited(err error) bool { return hasStatus(err, 429) }

// IsServerError reports whether err is an APIError with a 5xx status.
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status >= 500
	}
	return false
}

// IsNetworkError reports whether err is a NetworkError (request never completed).
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

func hasStatus(err error, status int) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status == status
	}
	return false
}

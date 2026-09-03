package noxaeapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RetryOptions controls retry behavior for network errors, 429s, and 5xx responses.
type RetryOptions struct {
	// Attempts is the max number of attempts including the first one. Default 3.
	Attempts int
	// BaseDelay is the base delay used for exponential backoff. Default 300ms.
	BaseDelay time.Duration
	// MaxDelay is the upper bound for any single backoff delay. Default 5s.
	MaxDelay time.Duration
}

// DefaultRetryOptions returns the SDK's default retry behavior.
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{Attempts: 3, BaseDelay: 300 * time.Millisecond, MaxDelay: 5 * time.Second}
}

// ClientOptions configures a Client or NetworkHubClient.
type ClientOptions struct {
	// BaseURL of the server, e.g. "http://localhost:8080". Required.
	BaseURL string
	// APIKey configured on the server (sent as the "key" header).
	APIKey string
	// Timeout is the per-request timeout. Default 10s.
	Timeout time.Duration
	// Retry controls retry behavior. Nil disables retries; use
	// DefaultRetryOptions() for the SDK's defaults.
	Retry *RetryOptions
	// Headers are extra headers sent on every request.
	Headers map[string]string
	// HTTPClient overrides the underlying *http.Client, mainly for testing
	// (e.g. a custom Transport with a mock RoundTripper).
	HTTPClient *http.Client
}

type httpEngine struct {
	baseURL      string
	apiKey       string
	extraHeaders map[string]string
	retry        *RetryOptions
	client       *http.Client
}

func newHTTPEngine(opts ClientOptions) (*httpEngine, error) {
	if opts.BaseURL == "" {
		return nil, &ConfigError{Message: "NoxAeApiClient requires a non-empty BaseURL"}
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	return &httpEngine{
		baseURL:      strings.TrimRight(opts.BaseURL, "/"),
		apiKey:       opts.APIKey,
		extraHeaders: opts.Headers,
		retry:        opts.Retry,
		client:       client,
	}, nil
}

// ConfigError indicates invalid client configuration (e.g. empty BaseURL).
type ConfigError struct{ Message string }

func (e *ConfigError) Error() string { return "noxaeapi: config error: " + e.Message }

type bodyKind int

const (
	bodyNone bodyKind = iota
	bodyForm          // application/x-www-form-urlencoded (Javalin's ctx.formParam(...))
	bodyJSON          // application/json (Javalin's ctx.bodyAsClass(...))
)

type requestBody struct {
	kind bodyKind
	form url.Values
	json interface{}
}

func noBody() requestBody { return requestBody{kind: bodyNone} }

func formBody(values url.Values) requestBody {
	return requestBody{kind: bodyForm, form: values}
}

func jsonBody(v interface{}) requestBody {
	return requestBody{kind: bodyJSON, json: v}
}

// buildURL joins the base URL, the /v1 prefix, path, and query params.
func (h *httpEngine) buildURL(path string, query url.Values) string {
	u := h.baseURL + "/v1/" + strings.TrimLeft(path, "/")
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return u
}

// doJSON performs a request and unmarshals the JSON response into out.
// Pass a pointer, or nil if the response body doesn't matter.
func (h *httpEngine) doJSON(ctx context.Context, method, path string, body requestBody, query url.Values, out interface{}) error {
	raw, err := h.doRaw(ctx, method, path, body, query)
	if err != nil {
		return err
	}
	if out == nil || raw == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

// doRaw performs a request and returns the raw response body (nil for empty/204 responses).
func (h *httpEngine) doRaw(ctx context.Context, method, path string, body requestBody, query url.Values) ([]byte, error) {
	fullURL := h.buildURL(path, query)
	maxAttempts := 1
	if h.retry != nil {
		maxAttempts = h.retry.Attempts
		if maxAttempts < 1 {
			maxAttempts = 1
		}
	}

	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		var reqBody io.Reader
		contentType := ""

		switch body.kind {
		case bodyForm:
			encoded := body.form.Encode()
			reqBody = strings.NewReader(encoded)
			contentType = "application/x-www-form-urlencoded"
		case bodyJSON:
			data, err := json.Marshal(body.json)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewReader(data)
			contentType = "application/json"
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		for k, v := range h.extraHeaders {
			req.Header.Set(k, v)
		}
		if h.apiKey != "" {
			req.Header.Set("key", h.apiKey)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			netErr := &NetworkError{Method: method, Path: path, Err: err}
			if h.retry != nil && attempt < maxAttempts-1 {
				lastErr = netErr
				sleep(backoffDelay(attempt, *h.retry))
				continue
			}
			return nil, netErr
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			netErr := &NetworkError{Method: method, Path: path, Err: readErr}
			if h.retry != nil && attempt < maxAttempts-1 {
				lastErr = netErr
				sleep(backoffDelay(attempt, *h.retry))
				continue
			}
			return nil, netErr
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if resp.StatusCode == 204 || len(respBody) == 0 {
				return nil, nil
			}
			return respBody, nil
		}

		apiErr := &APIError{
			Status: resp.StatusCode,
			Method: method,
			Path:   path,
			Body:   respBody,
		}

		if resp.StatusCode == 429 {
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, err := strconv.ParseFloat(ra, 64); err == nil {
					ms := int64(secs * 1000)
					apiErr.RetryAfterMs = &ms
				}
			}
		}

		retryable := resp.StatusCode == 429 || resp.StatusCode >= 500
		if retryable && h.retry != nil && attempt < maxAttempts-1 {
			lastErr = apiErr
			var delay time.Duration
			if apiErr.RetryAfterMs != nil {
				delay = time.Duration(*apiErr.RetryAfterMs) * time.Millisecond
			} else {
				delay = backoffDelay(attempt, *h.retry)
			}
			sleep(delay)
			continue
		}

		return nil, apiErr
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &NetworkError{Method: method, Path: path, Err: io.ErrUnexpectedEOF}
}

func sleep(d time.Duration) {
	if d > 0 {
		time.Sleep(d)
	}
}

func backoffDelay(attempt int, opts RetryOptions) time.Duration {
	base := opts.BaseDelay
	if base <= 0 {
		base = 300 * time.Millisecond
	}
	maxDelay := opts.MaxDelay
	if maxDelay <= 0 {
		maxDelay = 5 * time.Second
	}
	exp := float64(base) * math.Pow(2, float64(attempt))
	if exp > float64(maxDelay) {
		exp = float64(maxDelay)
	}
	// full jitter
	return time.Duration(rand.Float64() * exp)
}

// values builds url.Values from key/value pairs, skipping empty-string values
// for optional params. Use for building query strings.
func values(pairs ...string) url.Values {
	v := url.Values{}
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			v.Set(pairs[i], pairs[i+1])
		}
	}
	return v
}

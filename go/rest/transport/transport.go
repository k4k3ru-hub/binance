// Package transport defines the contract used by Binance REST operations.
package transport

import (
	"context"
	"net/http"
	"net/url"
)

// Request is an immutable REST request value.
type Request struct {
	Method string
	Path   string
	Header http.Header
	Query  url.Values
}

// Executor executes REST requests.
type Executor interface {
	Do(ctx context.Context, request Request) ([]byte, error)
}

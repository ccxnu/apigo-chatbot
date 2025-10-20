package domain

import (
	"context"
)

type HTTPHeader struct {
	Key   string
	Value string
}

// HTTPRequest holds all data required to make an external HTTP request
type HTTPRequest struct {
	URL               string
	Method            string
	Body              any
	AuthKey           string
	AuthValue         string
	AdditionalHeaders []HTTPHeader
}

type HTTPClient interface {
	Do(ctx context.Context, req HTTPRequest, response any) error
}

package utils

import (
	"context"
	"net/http"
	"time"
)

// HTTPClient provides an interface for suplying the SDK with a custom HTTP client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}


type SDKConfig struct {
	Client            HTTPClient
	Security          func(context.Context) (interface{}, error)
	ServerURL         string
	ServerIndex       int
	ServerDefaults    []map[string]string
	Language          string
	OpenAPIDocVersion string
	SDKVersion        string
	GenVersion        string
	UserAgent         string
	// RetryConfig       *retry.Config
	// Hooks             *hooks.Hooks
	Timeout           *time.Duration
}

type WetroSDK struct {

}


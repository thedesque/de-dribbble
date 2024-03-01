package dribbble

import (
	"net/http"

	"golang.org/x/oauth2"
)

// Config for Dribbble Auth
type Config struct {
	Token      *oauth2.Token
	HTTPClient *http.Client
}

// NewConfig for auth
func NewConfig(token *oauth2.Token) *Config {
	return &Config{
		Token:      token,
		HTTPClient: http.DefaultClient,
	}
}

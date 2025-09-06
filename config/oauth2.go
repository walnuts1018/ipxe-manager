package config

import (
	"net/url"
)

type OAuth2Config struct {
	ClientID              string   `env:"CLIENT_ID,required"`
	ClientSecret          string   `env:"CLIENT_SECRET,required"`
	IntrospectionEndpoint *url.URL `env:"INTROSPECTION_ENDPOINT,required"`
	AllowedAudiences      []string `env:"ALLOWED_AUDIENCE" envSeparator:","`
}

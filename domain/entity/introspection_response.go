package entity

type IntrospectionResponse struct {
	Active     bool     `json:"active"`
	TokenType  string   `json:"token_type"`
	Issuer     string   `json:"iss"`
	Subject    string   `json:"sub"`
	Audience   []string `json:"aud"`
	Expiration int64    `json:"exp"`
	IssuedAt   int64    `json:"iat"`
	TokenID    string   `json:"jti"`
	ClientID   string   `json:"client_id"`
	Scope      string   `json:"scope"`
	Username   string   `json:"username"`
}

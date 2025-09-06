package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/walnuts1018/ipxe-manager/config"
	"github.com/walnuts1018/ipxe-manager/definitions"
	"github.com/walnuts1018/ipxe-manager/domain/entity"
	"github.com/walnuts1018/ipxe-manager/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type AuthService struct {
	ClientID                                  string
	ClientSecret                              string
	IntrospectionEndpoint                     *url.URL
	IntrospectionEndpointAuthMethodsSupported []string

	client *http.Client
}

func NewAuthService(ctx context.Context, cfg config.OAuth2Config) (*AuthService, error) {
	// ctx, span := tracer.Tracer.Start(ctx, "NewAuthService")
	// defer span.End()

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.JoinPath("/.well-known/oauth-authorization-server").String(), nil)
	// if err != nil {
	// 	return nil, err
	// }
	// req.Header.Set("Accept", "application/json")
	// req.Header.Set("User-Agent", definitions.UserAgent)

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer func() {
	// 	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	// 	resp.Body.Close()              //nolint:errcheck
	// }()

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("failed to get issuer info: %s", resp.Status)
	// }

	// var issuerInfo struct {
	// 	IntrospectionEndpoint                     string   `json:"introspection_endpoint"`
	// 	IntrospectionEndpointAuthMethodsSupported []string `json:"introspection_endpoint_auth_methods_supported"`
	// }
	// if err := json.NewDecoder(resp.Body).Decode(&issuerInfo); err != nil {
	// 	return nil, err
	// }

	return &AuthService{
		ClientID:              cfg.ClientID,
		ClientSecret:          cfg.ClientSecret,
		IntrospectionEndpoint: cfg.IntrospectionEndpoint,
		IntrospectionEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post"},
		client: client,
	}, nil
}

var ErrUnknownAuthMethod = errors.New("unknown auth method for introspection endpoint")

func (s *AuthService) IntrospectToken(ctx context.Context, token string) (entity.IntrospectionResponse, error) {
	switch {
	case slices.Contains(s.IntrospectionEndpointAuthMethodsSupported, "client_secret_basic"):
		return s.introspectTokenWithClientSecretBasic(ctx, token)
	case slices.Contains(s.IntrospectionEndpointAuthMethodsSupported, "client_secret_post"):
		return s.introspectTokenWithClientSecretPost(ctx, token)
	default:
		return entity.IntrospectionResponse{}, fmt.Errorf("%w: %v", ErrUnknownAuthMethod, s.IntrospectionEndpointAuthMethodsSupported)
	}
}

func (s *AuthService) introspectTokenWithClientSecretBasic(ctx context.Context, accessToken string) (entity.IntrospectionResponse, error) {
	data := url.Values{}
	data.Set("token", accessToken)
	data.Set("token_type_hint", "access_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.IntrospectionEndpoint.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return entity.IntrospectionResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.ClientID, s.ClientSecret)))))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", definitions.UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return entity.IntrospectionResponse{}, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()              //nolint:errcheck
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to read response body", slog.String("error", err.Error()))
			return entity.IntrospectionResponse{}, fmt.Errorf("failed to introspect token: %s", resp.Status)
		}
		slog.Debug("response status", slog.String("response", string(body)), slog.String("status", resp.Status))
		return entity.IntrospectionResponse{}, fmt.Errorf("failed to introspect token: %s", resp.Status)
	}

	var introspectionResponse entity.IntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspectionResponse); err != nil {
		return entity.IntrospectionResponse{}, err
	}

	return introspectionResponse, nil
}

func (s *AuthService) introspectTokenWithClientSecretPost(ctx context.Context, accessToken string) (entity.IntrospectionResponse, error) {
	data := url.Values{}
	data.Set("token", accessToken)
	data.Set("token_type_hint", "access_token")
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.IntrospectionEndpoint.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return entity.IntrospectionResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", definitions.UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return entity.IntrospectionResponse{}, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()              //nolint:errcheck
	}()

	if resp.StatusCode != http.StatusOK {
		return entity.IntrospectionResponse{}, fmt.Errorf("failed to introspect token: %s", resp.Status)
	}

	var introspectionResponse entity.IntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspectionResponse); err != nil {
		return entity.IntrospectionResponse{}, err
	}

	return introspectionResponse, nil
}

var _ usecase.AuthService = (*AuthService)(nil)

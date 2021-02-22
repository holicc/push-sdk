package http

import (
	"context"
)

type AuthClient struct {
	client *HTTPClient
	token  Token
}

type Token interface {
	TokenRequest() ([]byte, error)
	ParseResponse([]byte) error

	SetHeader() []HTTPOption

	GetAuthUrl() string
	GetAuthMethod() string
	GetAccessToken() string

	IsValidate() bool
}

func NewAuthClient(t Token) *AuthClient {
	return &AuthClient{
		token:  t,
		client: NewHTTPClient(),
	}
}

func (ac *AuthClient) GetAuthToken(ctx context.Context) (string, error) {
	if ac.token.IsValidate() {
		return ac.token.GetAccessToken(), nil
	}

	body, err := ac.token.TokenRequest()
	if err != nil {
		return "", err
	}

	request := &PushRequest{
		Method: ac.token.GetAuthMethod(),
		URL:    ac.token.GetAuthUrl(),
		Body:   body,
		Header: ac.token.SetHeader(),
	}

	resp, err := ac.client.Do(ctx, request)
	if err != nil {
		return "", err
	}

	if resp.Status == 200 {
		err = ac.token.ParseResponse(resp.Body)
		if err != nil {
			return "", err
		}
		return ac.token.GetAccessToken(), nil
	}
	return "", nil
}

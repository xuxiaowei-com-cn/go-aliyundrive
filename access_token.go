package aliyundrive

import (
	"net/http"
)

type AccessTokenBody struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AccessToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type AccessTokenService struct {
	client *Client
}

// PostAccessToken 授权 code 获取 access_token https://www.yuque.com/aliyundrive/zpfszx/efabcs#Fyis9
func (s *AccessTokenService) PostAccessToken(body AccessTokenBody) (*AccessToken, *Response, error) {

	u := "oauth/access_token"

	req, err := s.client.NewRequest(http.MethodPost, u, nil, body)
	if err != nil {
		return nil, nil, err
	}

	var accessToken *AccessToken
	resp, err := s.client.Do(req, &accessToken)
	if err != nil {
		return nil, resp, err
	}

	return accessToken, resp, nil
}

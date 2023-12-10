package aliyundrive

import (
	"fmt"
	"net/http"
)

type QrCodeBody struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"` // user:base、file:all:read、file:all:write
	Width        int      `json:"width"`
	Height       int      `json:"height"`
}

type QrCode struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	Sid       string `json:"sid"`
}

type QrCodeService struct {
	client *Client
}

// PostQrCode 获取授权二维码 https://www.yuque.com/aliyundrive/zpfszx/ttfoy0xt2pza8lof#XuClO
func (s *QrCodeService) PostQrCode(body QrCodeBody) (*QrCode, *Response, error) {

	u := "oauth/authorize/qrcode"

	req, err := s.client.NewRequest(http.MethodPost, u, nil, body)
	if err != nil {
		return nil, nil, err
	}

	var qrCode *QrCode
	resp, err := s.client.Do(req, &qrCode)
	if err != nil {
		return nil, resp, err
	}

	return qrCode, resp, nil
}

type QrCodeStatus struct {
	Status   string `json:"status"`
	AuthCode string `json:"authCode"`
}

// GetQrCodeStatus 获取二维码登录状态 https://www.yuque.com/aliyundrive/zpfszx/ttfoy0xt2pza8lof#MW79B
func (s *QrCodeService) GetQrCodeStatus(sid string) (*QrCodeStatus, *Response, error) {

	u := fmt.Sprintf("oauth/qrcode/%s/status", sid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var qrCodeStatus *QrCodeStatus
	resp, err := s.client.Do(req, &qrCodeStatus)
	if err != nil {
		return nil, resp, err
	}

	return qrCodeStatus, resp, nil
}

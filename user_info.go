package aliyundrive

import (
	"net/http"
)

type UsersInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Phone  string `json:"phone"`
}

type UsersInfoService struct {
	client *Client
}

// GetUsersInfo 获取用户信息 https://www.yuque.com/aliyundrive/zpfszx/mbb50w#xZ6HQ
func (s *UsersInfoService) GetUsersInfo() (*UsersInfo, *Response, error) {

	u := "oauth/users/info"

	req, err := s.client.NewRequest(http.MethodGet, u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var usersInfo *UsersInfo
	resp, err := s.client.Do(req, &usersInfo)
	if err != nil {
		return nil, resp, err
	}

	return usersInfo, resp, nil
}

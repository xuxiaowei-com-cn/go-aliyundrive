package aliyundrive

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestPostAccessToken(t *testing.T) {

	var clientId = os.Getenv("GO_ALIYUNDRIVE_CLIENT_ID")
	var clientSecret = os.Getenv("GO_ALIYUNDRIVE_CLIENT_SECRET")
	var code = os.Getenv("GO_ALIYUNDRIVE_CODE")

	client, err := NewClient("")
	assert.NoError(t, err)

	body := AccessTokenBody{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    "authorization_code",
		Code:         code,
	}

	accessToken, response, err := client.AccessToken.PostAccessToken(body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	t.Log("AccessToken", accessToken.AccessToken)
	t.Log("RefreshToken", accessToken.RefreshToken)

	assert.NotEqual(t, "", accessToken.AccessToken)
	assert.NotEqual(t, "", accessToken.RefreshToken)

	body = AccessTokenBody{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    "refresh_token",
		RefreshToken: accessToken.RefreshToken,
	}

	accessToken, response, err = client.AccessToken.PostAccessToken(body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	t.Log("AccessToken", accessToken.AccessToken)
	t.Log("RefreshToken", accessToken.RefreshToken)

	assert.NotEqual(t, "", accessToken.AccessToken)
	assert.NotEqual(t, "", accessToken.RefreshToken)
}

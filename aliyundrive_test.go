package aliyundrive

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestGetQrCodeStatusPostAccessToken(t *testing.T) {

	var sid = os.Getenv("GO_ALIYUNDRIVE_SID")

	client, err := NewClient("")
	assert.NoError(t, err)

	qrCodeStatus, response, err := client.QrCode.GetQrCodeStatus(sid)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	t.Log("Status", qrCodeStatus.Status)
	t.Log("AuthCode", qrCodeStatus.AuthCode)

	assert.NotEqual(t, "", qrCodeStatus.Status)
	assert.NotEqual(t, "", qrCodeStatus.AuthCode)

	var clientId = os.Getenv("GO_ALIYUNDRIVE_CLIENT_ID")
	var clientSecret = os.Getenv("GO_ALIYUNDRIVE_CLIENT_SECRET")
	var grantType = "authorization_code"

	body := AccessTokenBody{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    grantType,
		Code:         qrCodeStatus.AuthCode,
	}

	accessToken, response, err := client.AccessToken.PostAccessToken(body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	t.Log("AccessToken", accessToken.AccessToken)
	t.Log("RefreshToken", accessToken.RefreshToken)

	assert.NotEqual(t, "", accessToken.AccessToken)
	assert.NotEqual(t, "", accessToken.RefreshToken)
}

package aliyundrive

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestPostQrCode(t *testing.T) {

	var clientId = os.Getenv("GO_ALIYUNDRIVE_CLIENT_ID")
	var clientSecret = os.Getenv("GO_ALIYUNDRIVE_CLIENT_SECRET")
	var scopes = strings.Split(os.Getenv("GO_ALIYUNDRIVE_CLIENT_SCOPES"), ",")
	var height = GetenvInt("GO_ALIYUNDRIVE_HEIGHT", 430)
	var width = GetenvInt("GO_ALIYUNDRIVE_WIDTH", 430)

	client, err := NewClient("")
	assert.NoError(t, err)

	body := QrCodeBody{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Height:       height,
		Width:        width,
	}

	qrCode, response, err := client.QrCode.PostQrCode(body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	t.Log("QrCodeUrl", qrCode.QrCodeUrl)
	t.Log("Sid", qrCode.Sid)

	assert.NotEqual(t, "", qrCode.QrCodeUrl)
	assert.NotEqual(t, "", qrCode.Sid)
}

func TestGetQrCodeStatus(t *testing.T) {

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
}

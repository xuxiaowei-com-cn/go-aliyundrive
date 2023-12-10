package aliyundrive

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestGetUsersInfo(t *testing.T) {

	var token = os.Getenv("GO_ALIYUNDRIVE_TOKEN")

	client, err := NewClient(token)
	assert.NoError(t, err)

	usersInfo, response, err := client.UsersInfo.GetUsersInfo()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	t.Log("Id", usersInfo.Id)
	t.Log("Name", usersInfo.Name)
	t.Log("Avatar", usersInfo.Avatar)
	t.Log("Phone", usersInfo.Phone)

	assert.NotEqual(t, "", usersInfo.Id)
	assert.NotEqual(t, "", usersInfo.Name)
	assert.NotEqual(t, "", usersInfo.Avatar)
}

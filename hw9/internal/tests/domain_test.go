package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	userSecond, err := client.createUser("user_2", "email_2@gmail.com")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(userSecond.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	userSecond, err := client.createUser("user_2", "email_second@gmail.com")
	assert.NoError(t, err)

	_, err = client.updateAd(userSecond.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

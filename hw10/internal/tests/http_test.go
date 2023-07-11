package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests for new http methods from homework 9 and 8 are added here

func TestGetAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.getAd(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.deleteAd(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "")
	assert.Equal(t, response.Data.Text, "")
}

func TestListAdsWithFilters(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(user.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsWithFilters(true, user.Data.ID, "")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestSearchAds(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello cats", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.searchAds("cats")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, response.Data.ID)
	assert.Equal(t, ads.Data[0].Title, response.Data.Title)
	assert.Equal(t, ads.Data[0].Text, response.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, response.Data.AuthorID)
}

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.NickName, "user_1")
	assert.Equal(t, response.Data.Email, "email@gmail.com")
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	response, err = client.getUser(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.NickName, "user_1")
	assert.Equal(t, response.Data.Email, "email@gmail.com")
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("user_1", "email@gmail.com")
	assert.NoError(t, err)

	response, err = client.deleteUser(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.NickName, "")
	assert.Equal(t, response.Data.Email, "")
}

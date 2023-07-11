package localrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/domain/models"
	"testing"
)

type TestSuite struct {
	suite.Suite
	user *models.User
}

func (suite *TestSuite) SetupTest() {
	suite.user = &models.User{
		ID:       0,
		NickName: "test nickname",
		Email:    "test email",
	}
}

func (suite *TestSuite) TestGetUser() {
	userRepo := NewUserRepo()
	userRepo.storage[0] = suite.user

	user, err := userRepo.GetUser(context.Background(), 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), *suite.user, *user)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

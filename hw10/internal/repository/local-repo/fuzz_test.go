package localrepo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func FuzzIntGetUser(f *testing.F) {
	userRepo := NewUserRepo()

	testcases := []int{90, 1000}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s int) {
		user, err := userRepo.GetUser(context.Background(), int64(s))

		assert.Nil(t, user)
		assert.Equal(t, fmt.Errorf("the user does not exist"), err)
	})
}

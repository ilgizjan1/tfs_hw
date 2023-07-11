package localrepo

import (
	"context"
	"homework10/internal/domain/models"
	"testing"
)

func Users() {
	repo := NewUserRepo()

	user := models.User{
		NickName: "testing",
		Email:    "testing@gmail.com",
	}

	for i := 0; i < 100; i++ {
		_, _ = repo.AddUser(context.Background(), user)
		_ = repo.Delete(context.Background(), 0)
	}
}

func UsersFast() {
	repo := NewUserRepo()

	user := models.User{
		NickName: "testing",
		Email:    "testing@gmail.com",
	}

	for i := 0; i < 100; i++ {
		repo.storage[0] = &user
		delete(repo.storage, 0)
	}
}

func BenchmarkUsers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Users()
	}
}

func BenchmarkUsersFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UsersFast()
	}
}

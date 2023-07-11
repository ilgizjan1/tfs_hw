package localrepo

import (
	"context"
	"fmt"
	"homework9/internal/domain/models"
	"sync"
)

type UserRepo struct {
	storage    map[int64]*models.User
	lastUserID int64
	mutex      sync.Mutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{storage: make(map[int64]*models.User), lastUserID: -1}
}

func (r *UserRepo) GetUser(ctx context.Context, id int64) (*models.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		user, ok := r.storage[id]
		if !ok {
			return nil, fmt.Errorf("the user does not exist")
		}
		return user, nil
	}
}

func (r *UserRepo) AddUser(ctx context.Context, user models.User) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.lastUserID++
		r.storage[r.lastUserID] = &user
		return r.lastUserID, nil
	}
}

func (r *UserRepo) Update(ctx context.Context, userID int64, nickName string, email string) (*models.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.storage[userID].NickName = nickName
		r.storage[userID].Email = email
		return r.storage[userID], nil
	}
}

func (r *UserRepo) Delete(ctx context.Context, userID int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		delete(r.storage, userID)
		return nil
	}
}

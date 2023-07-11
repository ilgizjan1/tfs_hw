package userRepo

import (
	"context"
	"fmt"
	"homework8/internal/users"
	"sync"
)

type RepositoryLocal struct {
	storage    map[int64]*users.User
	lastUserID int64
	mutex      sync.Mutex
}

func New() *RepositoryLocal {
	return &RepositoryLocal{storage: make(map[int64]*users.User), lastUserID: -1}
}

func (r *RepositoryLocal) GetUser(ctx context.Context, id int64) (*users.User, error) {
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

func (r *RepositoryLocal) AddUser(ctx context.Context, user users.User) (int64, error) {
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

func (r *RepositoryLocal) Update(ctx context.Context, userID int64, nickName string, email string) (*users.User, error) {
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

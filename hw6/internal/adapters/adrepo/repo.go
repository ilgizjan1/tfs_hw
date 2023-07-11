package adrepo

import (
	"context"
	"fmt"
	"homework6/internal/ads"
	"sync"
)

func New() *RepositoryLocal {
	return &RepositoryLocal{storage: make(map[int64]*ads.Ad), lastAdID: -1}
}

type RepositoryLocal struct {
	storage  map[int64]*ads.Ad
	lastAdID int64
	mutex    sync.Mutex
}

func (r *RepositoryLocal) AddAd(ctx context.Context, ad ads.Ad) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.lastAdID++
		r.storage[r.lastAdID] = &ad
		return r.lastAdID, nil
	}
}

func (r *RepositoryLocal) GetAuthorID(ctx context.Context, adID int64) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		if r.lastAdID < adID {
			return 0, fmt.Errorf("there is no such ad")
		}
		return r.storage[adID].AuthorID, nil
	}
}

func (r *RepositoryLocal) SetStatus(ctx context.Context, adID int64, published bool) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.storage[adID].Published = published
		return r.storage[adID], nil
	}
}

func (r *RepositoryLocal) Update(ctx context.Context, adID int64, title string, text string) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.storage[adID].Title = title
		r.storage[adID].Text = text
		return r.storage[adID], nil
	}
}

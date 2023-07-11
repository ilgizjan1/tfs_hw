package localrepo

import (
	"context"
	"fmt"
	"homework9/internal/domain/models"
	"sync"
	"time"
)

type AdRepo struct {
	storage  map[int64]*models.Ad
	lastAdID int64
	mutex    sync.Mutex
}

func NewAdRepo() *AdRepo {
	return &AdRepo{storage: make(map[int64]*models.Ad), lastAdID: -1}
}

func (r *AdRepo) GetAd(ctx context.Context, adID int64) (*models.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		ad, ok := r.storage[adID]
		if !ok {
			return nil, fmt.Errorf("the ad does not exist")
		}
		return ad, nil
	}
}

func (r *AdRepo) GetAds(ctx context.Context) ([]*models.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		adSlice := make([]*models.Ad, 0)
		r.mutex.Lock()
		defer r.mutex.Unlock()
		for _, ad := range r.storage {
			adSlice = append(adSlice, ad)
		}
		return adSlice, nil
	}
}

func (r *AdRepo) AddAd(ctx context.Context, ad models.Ad) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.lastAdID++
		r.storage[r.lastAdID] = &ad
		r.storage[r.lastAdID].ID = r.lastAdID
		return r.lastAdID, nil
	}
}

func (r *AdRepo) SetStatus(ctx context.Context, adID int64, published bool) (*models.Ad, error) {
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

func (r *AdRepo) Update(ctx context.Context, adID int64, title string, text string) (*models.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.storage[adID].Title = title
		r.storage[adID].Text = text
		r.storage[adID].DateUpdate = time.Now().UTC()
		return r.storage[adID], nil
	}
}

func (r *AdRepo) DeleteAd(ctx context.Context, adID int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		delete(r.storage, adID)
		return nil
	}
}

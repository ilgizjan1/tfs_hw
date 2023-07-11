package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"homework10/internal/domain/models"
	repoMock "homework10/internal/service/mock"
	"reflect"
	"testing"
	"time"
)

func TestErrNoAccess_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNoAccessAd to string",
			err:      ErrNoAccessAd,
			expected: fmt.Sprint(ErrNoAccessAd),
		},
		{
			name:     "empty error",
			err:      fmt.Errorf(""),
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrNoAccess{
				Err: tt.err,
			}
			assert.Equalf(t, tt.expected, e.Error(), "Error()")
		})
	}
}

func Test_checkPublished(t *testing.T) {
	tests := []struct {
		name      string
		published bool
		expected  funCheckAd
	}{
		{
			name:      "test for published = true",
			published: true,
			expected: func(ad models.Ad) bool {
				return ad.Published == true
			},
		},
		{
			name:      "test for published = false",
			published: false,
			expected: func(ad models.Ad) bool {
				return ad.Published == false
			},
		},
	}
	ad := models.Ad{Published: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPublished(tt.published); !reflect.DeepEqual(got(ad), tt.expected(ad)) {
				t.Errorf("checkPublished() = %v, inUser %v", got(ad), tt.expected(ad))
			}
		})
	}
}

func Test_checkUserID(t *testing.T) {
	tests := []struct {
		name     string
		userID   int64
		expected funCheckAd
	}{
		{
			name:   "test for userID = 0",
			userID: 0,
			expected: func(ad models.Ad) bool {
				return ad.UserID == 0
			},
		},
		{
			name:   "test for userID = 100",
			userID: 100,
			expected: func(ad models.Ad) bool {
				return ad.UserID == 100
			},
		},
	}
	ad := models.Ad{UserID: 0}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkUserID(tt.userID); !reflect.DeepEqual(got(ad), tt.expected(ad)) {
				t.Errorf("checkUserID() = %v, inUser %v", got(ad), tt.expected(ad))
			}
		})
	}
}

func Test_checkDate(t *testing.T) {
	now := time.Now()
	testTime, _ := time.Parse("2/1/2006", "31/7/2015")
	tests := []struct {
		name     string
		Date     time.Time
		expected funCheckAd
	}{
		{
			name: "test for true date",
			Date: now,
			expected: func(ad models.Ad) bool {
				return ad.DateCreation == now.Format(dateFormat)
			},
		},
		{
			name: "test for false date",
			Date: testTime,
			expected: func(ad models.Ad) bool {
				return ad.DateCreation == testTime.Format(dateFormat)
			},
		},
	}
	ad := models.Ad{UserID: 0}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDate(tt.Date); !reflect.DeepEqual(got(ad), tt.expected(ad)) {
				t.Errorf("checkDate() = %v, inUser %v", got(ad), tt.expected(ad))
			}
		})
	}
}

func Test_check(t *testing.T) {
	fun1 := checkPublished(true)
	tests := []struct {
		name      string
		ad        models.Ad
		functions []funCheckAd
		expected  bool
	}{
		{
			name:      "test for true",
			ad:        models.Ad{Published: true},
			functions: []funCheckAd{fun1},
			expected:  true,
		},
		{
			name:      "test for false",
			ad:        models.Ad{Published: false},
			functions: []funCheckAd{fun1},
			expected:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := check(tt.ad, tt.functions); got != tt.expected {
				t.Errorf("check() = %v, inUser %v", got, tt.expected)
			}
		})
	}
}

func TestCreateAd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)
	testTable := []struct {
		name    string
		inAd    models.Ad
		repoErr error
		wantErr bool
	}{
		{
			name: "true createAd",
			inAd: models.Ad{
				ID:           0,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       111,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "validate error",
			inAd: models.Ad{
				ID:           0,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       111,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			repoErr: fmt.Errorf("error from the data repository"),
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().AddAd(ctx, testCase.inAd).Return(testCase.inAd.ID, testCase.repoErr).Times(1)

			ad, err := adService.CreateAd(ctx, testCase.inAd.Title, testCase.inAd.Text, testCase.inAd.UserID)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ad)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.inAd, *ad)
		})
	}
}

func TestCreateAd_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)

	testTable := []struct {
		name string
		inAd models.Ad
	}{
		{
			name: "error in Title validation",
			inAd: models.Ad{
				ID:     10,
				Title:  "",
				Text:   "test Text",
				UserID: 111,
			},
		},
		{
			name: "error in Text validation",
			inAd: models.Ad{
				ID:     10,
				Title:  "test Title",
				Text:   "",
				UserID: 111,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			ad, err := adService.CreateAd(ctx, testCase.inAd.Title, testCase.inAd.Text, testCase.inAd.UserID)

			assert.Error(t, err)
			assert.Nil(t, ad)
		})
	}
}

func TestGetAd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)
	testTable := []struct {
		name    string
		adID    int64
		outAd   *models.Ad
		repoErr error
		wantErr bool
	}{
		{
			name: "true createAd",
			adID: 100,
			outAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       111,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name:    "validate error",
			adID:    100,
			outAd:   nil,
			repoErr: fmt.Errorf("error from the data repository"),
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAd(ctx, testCase.adID).Return(testCase.outAd, testCase.repoErr).Times(1)

			ad, err := adService.GetAdByID(ctx, testCase.adID)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ad)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.outAd, *ad)
		})
	}
}

func TestChangeAdStatus_After(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)
	testTable := []struct {
		name         string
		inAd         models.Ad
		outGetAd     *models.Ad
		outSetStatus *models.Ad
		repoErr      error
		wantErr      bool
	}{
		{
			name: "true test for changeAdStatus",
			inAd: models.Ad{
				ID:        100,
				Published: true,
				UserID:    47,
			},
			outGetAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       47,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "true test for changeAdStatus",
			inAd: models.Ad{
				ID:        100,
				Published: true,
				UserID:    47,
			},
			outGetAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       47,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			outSetStatus: nil,
			repoErr:      fmt.Errorf("error from the data repository: SetStatus()"),
			wantErr:      true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAd(ctx, testCase.inAd.ID).Return(testCase.outGetAd, nil).Times(1)
			adRepo.EXPECT().SetStatus(ctx, testCase.inAd.ID, testCase.inAd.Published).
				Return(testCase.outGetAd, testCase.repoErr).Times(1)

			ad, err := adService.ChangeAdStatus(ctx, testCase.inAd.ID, testCase.inAd.UserID, testCase.inAd.Published)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ad)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.outGetAd, *ad)
		})
	}
}

func TestChangeAdStatus_Before(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)
	testTable := []struct {
		name     string
		inAd     models.Ad
		outGetAd *models.Ad
		repoErr  error
	}{
		{
			name: "wrong userID",
			inAd: models.Ad{
				ID:        100,
				Published: true,
				UserID:    47,
			},
			outGetAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       48,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			repoErr: nil,
		},
		{
			name: "true test for changeAdStatus",
			inAd: models.Ad{
				ID:        100,
				Published: true,
				UserID:    47,
			},
			outGetAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       47,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			repoErr: fmt.Errorf("error from the data repository: SetStatus()"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAd(ctx, testCase.inAd.ID).Return(testCase.outGetAd, testCase.repoErr).Times(1)

			ad, err := adService.ChangeAdStatus(ctx, testCase.inAd.ID, testCase.inAd.UserID, testCase.inAd.Published)
			assert.Error(t, err)
			assert.Nil(t, ad)
		})
	}
}

func TestUpdateAd_Before(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)
	testTable := []struct {
		name     string
		inAd     models.Ad
		outGetAd *models.Ad
		errGetAd error
	}{
		{
			name: "wrong userID",
			inAd: models.Ad{
				ID:        100,
				Published: true,
				UserID:    47,
			},
			outGetAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       48,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			errGetAd: nil,
		},
		{
			name: "true test for changeAdStatus",
			inAd: models.Ad{
				ID:     100,
				Title:  "title test",
				Text:   "text test",
				UserID: 47,
			},
			outGetAd: nil,
			errGetAd: fmt.Errorf("error from the data repository: GetAd()"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAd(ctx, testCase.inAd.ID).Return(testCase.outGetAd, testCase.errGetAd).Times(1)

			ad, err := adService.UpdateAd(ctx, testCase.inAd.ID, testCase.inAd.UserID, testCase.inAd.Title, testCase.inAd.Text)
			assert.Error(t, err)
			assert.Nil(t, ad)
		})
	}
}

func TestUpdateAd_After(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)
	testTable := []struct {
		name          string
		inAd          models.Ad
		outAd         *models.Ad
		errRepoUpdate error
		wantErr       bool
	}{
		{
			name: "error from repository Update()",
			inAd: models.Ad{
				ID:        100,
				Published: true,
				UserID:    47,
			},
			outAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       47,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			errRepoUpdate: fmt.Errorf("error from repository Update()"),
			wantErr:       true,
		},
		{
			name: "error from Validate Title",
			inAd: models.Ad{
				ID:     100,
				Title:  "title test",
				Text:   "text test",
				UserID: 47,
			},
			outAd: &models.Ad{
				ID:           100,
				Title:        "",
				Text:         "test Text",
				UserID:       47,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			errRepoUpdate: nil,
			wantErr:       true,
		},
		{
			name: "true update",
			inAd: models.Ad{
				ID:     100,
				Title:  "title test",
				Text:   "text test",
				UserID: 47,
			},
			outAd: &models.Ad{
				ID:           100,
				Title:        "test Title",
				Text:         "test Text",
				UserID:       47,
				Published:    true,
				DateCreation: time.Now().UTC().Format(dateFormat),
				DateUpdate:   time.Now().UTC().Format(dateFormat),
			},
			errRepoUpdate: nil,
			wantErr:       false,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAd(ctx, testCase.inAd.ID).Return(testCase.outAd, nil).Times(1)

			adRepo.EXPECT().Update(ctx, testCase.outAd.ID, testCase.outAd.Title, testCase.outAd.Text).
				Return(testCase.outAd, testCase.errRepoUpdate).Times(1)

			ad, err := adService.UpdateAd(ctx, testCase.outAd.ID, testCase.outAd.UserID, testCase.outAd.Title, testCase.outAd.Text)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ad)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.outAd, *ad)
		})
	}
}

func TestDeleteAd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)

	testTable := []struct {
		name      string
		adID      int64
		userID    int64
		adGet     *models.Ad
		errGet    error
		errAccess error
		errDelete error
		wantErr   bool
	}{
		{
			name:   "true test delete",
			adID:   10,
			userID: 47,
			adGet: &models.Ad{
				ID:     10,
				UserID: 47,
			},
			wantErr: false,
		},
		{
			name:    "error test delete: GetAd()",
			adID:    10,
			userID:  47,
			adGet:   nil,
			errGet:  fmt.Errorf("error from repository GetAd()"),
			wantErr: true,
		},
		{
			name:   "error test delete: AccessUsers",
			adID:   10,
			userID: 100,
			adGet: &models.Ad{
				ID:     10,
				UserID: 47,
			},
			errAccess: ErrNoAccessAd,
			wantErr:   true,
		},
		{
			name:   "error test delete: DeleteAd()",
			adID:   10,
			userID: 47,
			adGet: &models.Ad{
				ID:     10,
				UserID: 47,
			},
			errDelete: fmt.Errorf("error from repository DeleteAd()"),
			wantErr:   true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAd(ctx, testCase.adID).Return(testCase.adGet, testCase.errGet).Times(1)

			if testCase.errGet == nil && testCase.errAccess == nil {
				adRepo.EXPECT().DeleteAd(ctx, testCase.adID).
					Return(testCase.errDelete).Times(1)
			}

			err := adService.DeleteAd(ctx, testCase.adID, testCase.userID)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

const (
	GetAd       = "GetAd"
	AccessUsers = "AccessUsers"
	DeleteAd    = "DeleteAd"
	GetAds      = "GetAds"
)

func TestGetAdsByTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)

	testTable := []struct {
		name         string
		search       string
		rawAds       []*models.Ad
		expected     []*models.Ad
		storageError map[string]error
		wantError    bool
	}{
		{
			name:   "true test",
			search: "cats",
			rawAds: []*models.Ad{
				{Title: "cats"},
				{Title: "cat"},
			},
			expected: []*models.Ad{
				{Title: "cats"},
			},
			storageError: make(map[string]error),
			wantError:    false,
		},
		{
			name:         "error from repository GetAds()",
			search:       "",
			rawAds:       nil,
			expected:     nil,
			storageError: map[string]error{GetAds: fmt.Errorf("error from repository GetAds()")},
			wantError:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAds(ctx).Return(testCase.rawAds, testCase.storageError[GetAds]).Times(1)

			ads, err := adService.GetAdsByTitle(ctx, testCase.search)
			if testCase.wantError {
				assert.Error(t, err)
				assert.Nil(t, ads)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.expected[0], *ads[0])
		})
	}
}

func TestListAds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adRepo := repoMock.NewMockAdRepository(ctrl)
	adService := NewAdService(adRepo)

	type filters struct {
		published    string
		userID       string
		dateCreation string
	}

	testTable := []struct {
		name         string
		filters      filters
		rawAds       []*models.Ad
		expected     []*models.Ad
		storageError map[string]error
		wantError    bool
	}{
		{
			name:         "error from repository GetAds()",
			rawAds:       nil,
			expected:     nil,
			storageError: map[string]error{GetAds: fmt.Errorf("error from repository GetAds()")},
			wantError:    true,
		},
		{
			name:    "error invalid filters: published",
			filters: filters{published: "abcdjsjsk"},
			rawAds: []*models.Ad{
				{Published: true}, {Published: false},
			},
			expected:  nil,
			wantError: true,
		},
		{
			name:    "error invalid filters: userID",
			filters: filters{userID: "absfsjf"},
			rawAds: []*models.Ad{
				{UserID: 10}, {UserID: 15},
			},
			expected:  nil,
			wantError: true,
		},
		{
			name:    "error invalid filters: date_creation",
			filters: filters{dateCreation: "absfsjf"},
			rawAds: []*models.Ad{
				{DateCreation: "absfsjf"}, {DateCreation: time.Now().Format(dateFormat)},
			},
			expected:  nil,
			wantError: true,
		},
		{
			name:    "true test: published - true",
			filters: filters{published: "true"},
			rawAds: []*models.Ad{
				{Published: true}, {Published: false},
			},
			expected: []*models.Ad{
				{Published: true},
			},
			wantError: false,
		},
		{
			name:    "true test: published - false",
			filters: filters{published: "false"},
			rawAds: []*models.Ad{
				{Published: true}, {Published: false},
			},
			expected: []*models.Ad{
				{Published: false},
			},
			wantError: false,
		},
		{
			name: "true test: empty filters => published = true",
			rawAds: []*models.Ad{
				{Published: true}, {Published: false},
			},
			expected: []*models.Ad{
				{Published: true},
			},
			wantError: false,
		},
		{
			name:    "true test: userID",
			filters: filters{userID: "10"},
			rawAds: []*models.Ad{
				{UserID: 10}, {UserID: 20},
			},
			expected: []*models.Ad{
				{UserID: 10},
			},
			wantError: false,
		},
		{
			name:    "true test: dateCreation",
			filters: filters{dateCreation: time.Now().Format(dateFormat)},
			rawAds: []*models.Ad{
				{DateCreation: time.Now().Format(dateFormat)}, {DateCreation: "31/7/2015"},
			},
			expected: []*models.Ad{
				{DateCreation: time.Now().Format(dateFormat)},
			},
			wantError: false,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			adRepo.EXPECT().GetAds(ctx).Return(testCase.rawAds, testCase.storageError[GetAds]).Times(1)

			ads, err := adService.ListAds(
				ctx,
				testCase.filters.published,
				testCase.filters.userID,
				testCase.filters.dateCreation,
			)

			if testCase.wantError {
				assert.Error(t, err)
				assert.Nil(t, ads)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.expected[0], *ads[0])
			assert.Equal(t, len(testCase.expected), len(ads))
		})
	}
}

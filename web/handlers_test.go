// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/repository"
	mockFilter "github.com/lvl484/positioning-filter/repository/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	log = logrus.New()
)

func newTestFilter(filterName, filterType string) *repository.Filter {
	conf, _ := json.Marshal("someString")
	userID, _ := uuid.Parse("d5cadefb-4d4d-4105-8244-1c354f936e69")

	return &repository.Filter{
		Name:          filterName,
		Type:          filterType,
		Configuration: conf,
		Reversed:      false,
		UserID:        userID,
	}
}

func TestHandlerAddFilterNoMock(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		log      *logrus.Logger
		wantCode int
	}{
		{
			name:     "TestAddFilterFailParseUUIDFromURL",
			log:      log,
			userID:   "1",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "TestAddFilterFailDecodeFilterFromBody",
			log:      log,
			userID:   "e5fadd54-7f09-11ea-bc55-0242ac130003",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			filters := mockFilter.NewMockFilters(ctrl)
			srv := httptest.NewServer(newRouter(filters, log))
			b, err := json.Marshal("1")
			assert.NoError(t, err)
			urlString := fmt.Sprintf("%s/users/%s/filters/", srv.URL, tt.userID)
			res, err := http.Post(urlString, "application/json", bytes.NewBuffer(b))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}

func TestHandlerAddFilterMock(t *testing.T) {
	tests := []struct {
		name      string
		filter    *repository.Filter
		log       *logrus.Logger
		mockError error
		wantCode  int
	}{
		{
			name:      "TestAddFilterFailDBCallAdd",
			log:       log,
			filter:    newTestFilter("Name", "rectangular"),
			mockError: errors.New("matrix error"),
			wantCode:  http.StatusInternalServerError,
		},
		{
			name:      "TestAddFilterSuccess",
			log:       log,
			filter:    newTestFilter("Paraska", "rectangular"),
			mockError: nil,
			wantCode:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			filters := mockFilter.NewMockFilters(ctrl)
			srv := httptest.NewServer(newRouter(filters, log))
			b, err := json.Marshal(tt.filter)
			assert.NoError(t, err)
			filters.EXPECT().Add(tt.filter).Return(tt.mockError)
			urlString := fmt.Sprintf("%s/users/%s/filters/", srv.URL, tt.filter.UserID)
			res, err := http.Post(urlString, "application/json", bytes.NewBuffer(b))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}

func TestHandlerGetOneFilterFailParseUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters, log))

	filter := newTestFilter("Name", "round")
	urlString := fmt.Sprintf("%s/users/%s/filters/%s", srv.URL, "err", filter.Name)
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestHandlerGetOneFilterMock(t *testing.T) {
	tests := []struct {
		name      string
		filter    *repository.Filter
		log       *logrus.Logger
		mockError error
		wantCode  int
	}{
		{
			name:      "TestGetOneFilterFailDBCallOneByUser",
			filter:    newTestFilter("Name", "rectangular"),
			log:       log,
			mockError: errors.New("Some Error"),
			wantCode:  http.StatusInternalServerError,
		},
		{
			name:      "TestGetOneFilterSuccess",
			filter:    newTestFilter("name", "rectangular"),
			log:       log,
			mockError: nil,
			wantCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			filters := mockFilter.NewMockFilters(ctrl)
			srv := httptest.NewServer(newRouter(filters, log))
			filters.EXPECT().OneByUser(tt.filter.UserID, tt.filter.Name).Return(tt.filter, tt.mockError)

			urlString := fmt.Sprintf("%s/users/%s/filters/%s", srv.URL, tt.filter.UserID, tt.filter.Name)
			res, err := http.Get(urlString)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}

func TestGetOffsetFailParseUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters, log))

	urlString := fmt.Sprintf("%s/users/%s/filters/", srv.URL, "err")
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestHandlerGetOffsetMock(t *testing.T) {
	filter1 := newTestFilter("Name1", "round")
	filter2 := newTestFilter("Name2", "round")
	both := []*repository.Filter{filter1, filter2}
	tests := []struct {
		name        string
		filter      *repository.Filter
		log         *logrus.Logger
		offset      int
		mockUserID  uuid.UUID
		mockError   error
		mockReturns []*repository.Filter
		wantCode    int
	}{
		{
			name:        "TestGetOffsetFailDBCallOneByUser",
			filter:      newTestFilter("name", "rectangular"),
			log:         log,
			offset:      0,
			mockUserID:  both[0].UserID,
			mockError:   errors.New("Some Error"),
			mockReturns: nil,
			wantCode:    http.StatusInternalServerError,
		},
		{
			name:        "TestGetOffsetSuccess",
			filter:      newTestFilter("name2", "rectangular"),
			log:         log,
			offset:      0,
			mockUserID:  both[0].UserID,
			mockError:   nil,
			mockReturns: both,
			wantCode:    http.StatusOK,
		},
		{
			name:        "TestGetOffsetFilterSuccess",
			filter:      newTestFilter("name3", "rectangular"),
			log:         log,
			offset:      -1512,
			mockUserID:  both[0].UserID,
			mockError:   nil,
			mockReturns: both,
			wantCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			filters := mockFilter.NewMockFilters(ctrl)
			srv := httptest.NewServer(newRouter(filters, log))

			filters.EXPECT().OffsetByUser(tt.mockUserID, 0).Return(tt.mockReturns, tt.mockError)

			urlString := fmt.Sprintf("%s/users/%s/filters/?offset=%v", srv.URL, tt.mockUserID, tt.offset)
			res, err := http.Get(urlString)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}

func TestHandlerUpdateFilterNoMock(t *testing.T) {
	tests := []struct {
		name     string
		filter   *repository.Filter
		log      *logrus.Logger
		userID   string
		wantCode int
	}{
		{
			name:     "TestUpdateFailParseUUIDFromURL",
			filter:   newTestFilter("name", "rectangular"),
			log:      log,
			userID:   "1",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "TestUpdateFailDecodeFilterFromBody",
			filter:   newTestFilter("name2", "rectangular"),
			log:      log,
			userID:   "e5fadd54-7f09-11ea-bc55-0242ac130003",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			filters := mockFilter.NewMockFilters(ctrl)
			srv := httptest.NewServer(newRouter(filters, log))

			b, err := json.Marshal("test")
			assert.NoError(t, err)

			urlString := fmt.Sprintf("%s/users/%v/filters/name", srv.URL, tt.userID)
			req, err := http.NewRequest(http.MethodPatch, urlString, bytes.NewBuffer(b))
			assert.NoError(t, err)

			client := &http.Client{}
			res, err := client.Do(req)
			assert.NoError(t, err)

			defer res.Body.Close()

			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}

func TestHandlerUpdateFilterMock(t *testing.T) {
	tests := []struct {
		name      string
		filter    *repository.Filter
		log       *logrus.Logger
		mockError error
		wantCode  int
	}{
		{
			name:      "TestUpdateFilterFailDBCallOneByUser",
			filter:    newTestFilter("name", "rectangular"),
			log:       log,
			mockError: errors.New("Some Error"),
			wantCode:  http.StatusInternalServerError,
		},
		{
			name:      "TestUpdateFilterSuccess",
			filter:    newTestFilter("name2", "rectangular"),
			log:       log,
			mockError: nil,
			wantCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			filters := mockFilter.NewMockFilters(ctrl)
			srv := httptest.NewServer(newRouter(filters, log))
			filters.EXPECT().Update(tt.filter).Return(tt.mockError)

			b, err := json.Marshal(&tt.filter)
			assert.NoError(t, err)

			urlString := fmt.Sprintf("%s/users/%s/filters/%s", srv.URL, tt.filter.UserID, tt.filter.Name)
			req, err := http.NewRequest(http.MethodPatch, urlString, bytes.NewBuffer(b))
			assert.NoError(t, err)
			client := &http.Client{}
			res, err := client.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}

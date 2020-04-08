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
	"github.com/stretchr/testify/assert"
)

func TestAddFilterSuccees(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	filters.EXPECT().Add(filter).Return(nil)

	b, err := json.Marshal(&filter)
	assert.NoError(t, err)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/", srv.URL, filter.UserID)
	res, err := http.Post(urlString, "application/json", bytes.NewBuffer(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestAddFilterFailDecode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	b, err := json.Marshal("b")
	assert.NoError(t, err)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/", srv.URL, filter.UserID)
	res, err := http.Post(urlString, "application/json", bytes.NewBuffer(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestAddFilterFailDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	filters.EXPECT().Add(filter).Return(errors.New("Error"))

	b, err := json.Marshal(&filter)
	assert.NoError(t, err)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/", srv.URL, filter.UserID)
	res, err := http.Post(urlString, "application/json", bytes.NewBuffer(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetOneFilterByUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	filters.EXPECT().OneByUser(filter.UserID, filter.Name).Return(filter, nil)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/%s", srv.URL, filter.UserID, filter.Name)
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetOneFilterByUserFailDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	filters.EXPECT().OneByUser(filter.UserID, filter.Name).Return(nil, errors.New("Error"))

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/%s", srv.URL, filter.UserID, filter.Name)
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetOneFilterByUserFailParseUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	urlString := fmt.Sprintf("%s/v1/router/%s/filters/%s", srv.URL, "err", filter.Name)
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestGetOffsetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter1 := newTestFilter("Name1", "round")
	filter2 := newTestFilter("Name2", "round")
	both := []*repository.Filter{filter1, filter2}

	filters.EXPECT().OffsetByUser(filter1.UserID, 0).Return(both, nil)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/", srv.URL, filter1.UserID)
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetOffsetFailDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter1 := newTestFilter("Name1", "round")

	filters.EXPECT().OffsetByUser(filter1.UserID, 0).Return(nil, errors.New("Err"))

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/", srv.URL, filter1.UserID)
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetOffsetFailParseUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/", srv.URL, "err")
	res, err := http.Get(urlString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestUpdateFilterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	filters.EXPECT().Update(filter).Return(nil)

	b, err := json.Marshal(&filter)
	assert.NoError(t, err)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/%s", srv.URL, filter.UserID, filter.Name)
	req, err := http.NewRequest(http.MethodPatch, urlString, bytes.NewBuffer(b))
	assert.NoError(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateFilterFailDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filters := mockFilter.NewMockFilters(ctrl)
	srv := httptest.NewServer(newRouter(filters))
	defer srv.Close()

	filter := newTestFilter("Name1", "round")
	filters.EXPECT().Update(filter).Return(errors.New("Err"))

	b, err := json.Marshal(&filter)
	assert.NoError(t, err)

	urlString := fmt.Sprintf("%s/v1/router/%s/filters/%s", srv.URL, filter.UserID, filter.Name)
	req, err := http.NewRequest(http.MethodPatch, urlString, bytes.NewBuffer(b))
	assert.NoError(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

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

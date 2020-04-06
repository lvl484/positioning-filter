// Package matcher provides matching positioning messages to a set of filtering rules
// and pereating messages to a topic in case of they are matched.
package matcher

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/position"
	"github.com/lvl484/positioning-filter/repository"
	mockRepository "github.com/lvl484/positioning-filter/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestXor(t *testing.T) {
	assert.False(t, xor(false, false))
	assert.False(t, xor(true, true))
	assert.True(t, xor(false, true))
	assert.True(t, xor(true, false))
}

func TestMatchRoundCriticalPosition(t *testing.T) {
	filter := newTestRoundFilter(-179, 0, 100, false)
	position := newTestPosition(178, 0)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundCriticalPosition2(t *testing.T) {
	filter := newTestRoundFilter(179, 0, 100, false)
	position := newTestPosition(-178, 0)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundCriticalPosition3(t *testing.T) {
	filter := newTestRoundFilter(0, 178, 100, false)
	position := newTestPosition(0, -178)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundCriticalPosition4(t *testing.T) {
	filter := newTestRoundFilter(0, -178, 100, false)
	position := newTestPosition(0, 178)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundCriticalPosition5(t *testing.T) {
	filter := newTestRoundFilter(-178, -178, 100, false)
	position := newTestPosition(178, 178)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundCriticalPosition6(t *testing.T) {
	filter := newTestRoundFilter(-179, 0, 10, false)
	position := newTestPosition(-170, 0)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundMatched(t *testing.T) {
	filter := newTestRoundFilter(0, 0, 50, false)
	position := newTestPosition(30, 40)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundNotMatched(t *testing.T) {
	filter := newTestRoundFilter(1, 1, 50, false)
	position := newTestPosition(-31.001, 41)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.False(t, matched)
}

func TestMatchRoundMatchedReversed(t *testing.T) {
	filter := newTestRoundFilter(0, 0, 50, true)
	position := newTestPosition(-30.001, 40)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRoundNotMatchedReversed(t *testing.T) {
	filter := newTestRoundFilter(0, 0, 130, true)
	position := newTestPosition(-50, 120)
	matched, err := matchRound(position, filter)
	assert.NoError(t, err)
	assert.False(t, matched)
}

func TestMatchRoundFail(t *testing.T) {
	filter := &repository.Filter{
		Name:          "Filter",
		Type:          "Round",
		Configuration: []byte{15, 13, 26, 77},
		Reversed:      true,
	}
	position := newTestPosition(-31, 40)
	matched, err := matchRound(position, filter)
	assert.NotNil(t, err)
	assert.False(t, matched)
}

func TestMatchRectangularWithConflictMatched(t *testing.T) {
	filter := newTestRectangularFilter(177, 4, -177, 1, false)
	position := newTestPosition(178, 2)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestCheckRectangularConflict(t *testing.T) {
	assert.True(t, isConflict(2, 1))
	assert.False(t, isConflict(1, 2))
}

func TestCheckRectangularConflictMatchedCriticalPosition(t *testing.T) {
	filter := newTestRectangularFilter(177, 4, -177, 1, false)
	position := newTestPosition(-178, 2)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestCheckRectangularConflictNotMatchedCriticalPosition(t *testing.T) {
	filter := newTestRectangularFilter(177, 4, -177, 1, false)
	position := newTestPosition(-176, 2)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.False(t, matched)
}

func TestMatchRectangularMatched(t *testing.T) {
	filter := newTestRectangularFilter(1, 4, 4, 1, false)
	position := newTestPosition(2, 2)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRectangularNotMatched(t *testing.T) {
	filter := newTestRectangularFilter(10, 15, -23, -48, false)
	position := newTestPosition(11, -49)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.False(t, matched)
}

func TestMatchRectangularMatchedReversed(t *testing.T) {
	filter := newTestRectangularFilter(0, 1, 1, 0, true)
	position := newTestPosition(14, 88)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatchRectangularNotMatchedReversed(t *testing.T) {
	filter := newTestRectangularFilter(0, 1, 1, 0, true)
	position := newTestPosition(0.1, 0.1)
	matched, err := matchRectangular(position, filter)
	assert.NoError(t, err)
	assert.False(t, matched)
}

func TestMatchRectangularFail(t *testing.T) {
	filter := &repository.Filter{
		Name:          "Filter",
		Type:          "Rectangular",
		Configuration: []byte{15, 13, 26, 77},
		Reversed:      true,
	}
	position := newTestPosition(-31, 40)
	matched, err := matchRectangular(position, filter)
	assert.NotNil(t, err)
	assert.False(t, matched)
}

func TestMatcherByTypeRound(t *testing.T) {
	matcher, err := matcherByType("round")
	assert.NoError(t, err)
	assert.NotNil(t, matcher)
}

func TestMatcherByTypeRectangular(t *testing.T) {
	matcher, err := matcherByType("rectangular")
	assert.NoError(t, err)
	assert.NotNil(t, matcher)
}

func TestMatcherByTypeFail(t *testing.T) {
	matcher, err := matcherByType("type")
	assert.EqualError(t, err, ErrBadFilterType)
	assert.Nil(t, matcher)
}

func TestMatcherMatchMatched(t *testing.T) {
	ctrl, filters := newTestFilterMock(t)
	defer ctrl.Finish()

	position := newTestPosition(14, 15)
	matcher := matcherFilters{filters: filters}

	filter1 := newTestRoundFilter(0, 0, 1, false)
	filter2 := newTestRoundFilter(0, 0, 100, false)
	filterSlice := []*repository.Filter{filter1, filter2}

	filters.EXPECT().AllByUser(position.UserID).Return(filterSlice, nil)
	matched, err := matcher.Match(position)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestMatcherMatchNotMatched(t *testing.T) {
	ctrl, filters := newTestFilterMock(t)
	defer ctrl.Finish()

	position := newTestPosition(10, 15)
	matcher := matcherFilters{filters: filters}

	filter1 := newTestRoundFilter(0, 0, 1, false)
	filter2 := newTestRoundFilter(0, 0, 100, true)
	filterSlice := []*repository.Filter{filter1, filter2}

	filters.EXPECT().AllByUser(position.UserID).Return(filterSlice, nil)
	matched, err := matcher.Match(position)
	assert.NoError(t, err)
	assert.False(t, matched)
}

func TestMatcherMatchFailAllByUser(t *testing.T) {
	ctrl, filters := newTestFilterMock(t)
	defer ctrl.Finish()

	position := newTestPosition(10, 15)
	matcher := matcherFilters{filters: filters}
	testError := errors.New("AllByUser failed because user died")

	filters.EXPECT().AllByUser(position.UserID).Return(nil, testError)
	matched, err := matcher.Match(position)
	assert.EqualError(t, err, "AllByUser failed because user died")
	assert.False(t, matched)
}

func TestMatcherMatchFailMatcherByType(t *testing.T) {
	ctrl, filters := newTestFilterMock(t)
	defer ctrl.Finish()

	position := newTestPosition(10, 15)
	matcher := matcherFilters{filters: filters}

	filter := newTestRoundFilter(0, 0, 100, true)
	filter.Type = "SomeType"
	filterSlice := []*repository.Filter{filter}

	filters.EXPECT().AllByUser(position.UserID).Return(filterSlice, nil)
	matched, err := matcher.Match(position)
	assert.EqualError(t, err, ErrBadFilterType)
	assert.False(t, matched)
}

func TestMatcherMatchFailMatch(t *testing.T) {
	ctrl, filters := newTestFilterMock(t)
	defer ctrl.Finish()

	position := newTestPosition(0, 0)
	matcher := matcherFilters{filters}

	filter := &repository.Filter{
		Name:          "Filter",
		Type:          "round",
		Configuration: []byte{15, 13, 26, 77},
		Reversed:      true,
	}
	filterSlice := []*repository.Filter{filter}

	filters.EXPECT().AllByUser(position.UserID).Return(filterSlice, nil)
	matched, err := matcher.Match(position)
	assert.NotNil(t, err)
	assert.False(t, matched)
}

func TestNewMatcher(t *testing.T) {
	ctrl, filters := newTestFilterMock(t)
	defer ctrl.Finish()

	matcher := NewMatcher(filters)
	assert.NotNil(t, matcher)
}

func newTestRoundFilter(centerLatitude, centerLongitude, radius float32, reversed bool) *repository.Filter {
	f := repository.RoundFilter{
		CenterLatitude:  centerLatitude,
		CentreLongitude: centerLongitude,
		Radius:          radius,
	}
	b, _ := json.Marshal(f)

	return &repository.Filter{
		Name:          "RoundFilter",
		Type:          "round",
		Configuration: b,
		Reversed:      reversed,
		UserID:        uuid.New(),
	}
}

func newTestRectangularFilter(
	topLeftLatitude, topLeftLongitude, bottomRightLatitude, bottomRightLongitude float32, reversed bool,
) *repository.Filter {
	f := repository.RectangularFilter{
		TopLeftLatitude:      topLeftLatitude,
		TopLeftLongitude:     topLeftLongitude,
		BottomRightLatitude:  bottomRightLatitude,
		BottomRightLongitude: bottomRightLongitude,
	}
	b, _ := json.Marshal(f)

	return &repository.Filter{
		Name:          "RectangularFilter",
		Type:          "rectangular",
		Configuration: b,
		Reversed:      reversed,
		UserID:        uuid.New(),
	}
}

func newTestPosition(latitude, longitude float32) position.Position {
	return position.Position{
		UserID:    uuid.New(),
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Millisecond),
	}
}

func newTestFilterMock(t *testing.T) (*gomock.Controller, *mockRepository.MockFilters) {
	ctrl := gomock.NewController(t)
	filters := mockRepository.NewMockFilters(ctrl)

	return ctrl, filters
}

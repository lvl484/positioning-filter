// Package matcher provides matching positioning messages to a set of filtering rules
// and pereating messages to a topic in case of they are matched.
package matcher

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/position"
	"github.com/lvl484/positioning-filter/repository"
	"github.com/stretchr/testify/assert"
)

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
	matcher := matcherByType("round")
	assert.NotNil(t, matcher)
}

func TestMatcherByTypeRectangular(t *testing.T) {
	matcher := matcherByType("rectangular")
	assert.NotNil(t, matcher)
}

func TestMatcherByTypeFail(t *testing.T) {
	matcher := matcherByType("type")
	assert.Nil(t, matcher)
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
		Type:          "Round",
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
		Type:          "Rectangular",
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

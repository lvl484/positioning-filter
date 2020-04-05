// Package matcher provides matching positioning messages to a set of filtering rules
// and pereating messages to a topic in case of they are matched.
package matcher

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/lvl484/positioning-filter/position"
	"github.com/lvl484/positioning-filter/repository"
)

const (
	ErrBadFilterType = "Bad type of filter"

	criticalLeftLatitude  float32 = -180
	criticalRightLatitude float32 = 180
)

type matcher func(position.Position, *repository.Filter) (bool, error)

type matcherFilters struct {
	filters repository.Filters
}

// Match checks if given position is matched with at least one filter
func (m matcherFilters) Match(pos position.Position) (bool, error) {
	filters, err := m.filters.AllByUser(pos.UserID)
	if err != nil {
		return false, err
	}

	for _, filter := range filters {
		match, err := matcherByType(filter.Type)
		if err != nil {
			return false, err
		}

		matched, err := match(pos, filter)

		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

// NewMatcher returns struct that implement Matcher interface
func NewMatcher(filters repository.Filters) Matcher {
	return matcherFilters{filters: filters}
}

func matcherByType(matcherType string) (matcher, error) {
	switch matcherType {
	case "round":
		return matchRound, nil
	case "rectangular":
		return matchRectangular, nil
	default:
		return nil, errors.New(ErrBadFilterType)
	}
}

func matchRectangular(pos position.Position, filter *repository.Filter) (bool, error) {
	var rfilter repository.RectangularFilter
	if err := json.Unmarshal(filter.Configuration, &rfilter); err != nil {
		return false, err
	}

	if isConflict(rfilter.TopLeftLatitude, rfilter.BottomRightLatitude) {
		delta := moveRectangularFilter(&rfilter)
		movePosition(&pos, delta)
	}

	matched := rfilter.BottomRightLatitude > pos.Latitude &&
		rfilter.TopLeftLatitude < pos.Latitude &&
		rfilter.BottomRightLongitude < pos.Longitude &&
		rfilter.TopLeftLongitude > pos.Longitude

	return xor(matched, filter.Reversed), nil
}

func matchRound(pos position.Position, filter *repository.Filter) (bool, error) {
	var rfilter repository.RoundFilter
	var limit float64 = 360
	if err := json.Unmarshal(filter.Configuration, &rfilter); err != nil {
		return false, err
	}
	switch {
	//round filter longitude and lalilude in critical position(near to +-180)
	case (limit/2-math.Abs(float64(rfilter.CenterLatitude))-float64(rfilter.Radius)) <= 0 && (limit/2-math.Abs(float64(rfilter.CentreLongitude))-float64(rfilter.Radius)) <= 0:
		matched := (limit-math.Abs(float64(pos.Latitude))-math.Abs(float64(rfilter.CenterLatitude)))*(limit-math.Abs(float64(pos.Latitude))-math.Abs(float64(rfilter.CenterLatitude)))+
			(limit-math.Abs(float64(pos.Longitude))-math.Abs(float64(rfilter.CentreLongitude)))*(limit-math.Abs(float64(pos.Longitude))-math.Abs(float64(rfilter.CentreLongitude))) <=
			float64(rfilter.Radius*rfilter.Radius)
		return xor(matched, filter.Reversed), nil
	//round filter lalilude in critical position(near to +-180)
	case (limit/2-math.Abs(float64(rfilter.CenterLatitude))-float64(rfilter.Radius)) <= 0 && !((limit/2 - math.Abs(float64(rfilter.CentreLongitude)) - float64(rfilter.Radius)) <= 0):
		matched := (limit-math.Abs(float64(pos.Latitude))-math.Abs(float64(rfilter.CenterLatitude)))*(limit-math.Abs(float64(pos.Latitude))-math.Abs(float64(rfilter.CenterLatitude)))+
			float64((pos.Longitude-rfilter.CentreLongitude)*(pos.Longitude-rfilter.CentreLongitude)) <=
			float64(rfilter.Radius*rfilter.Radius)
		return xor(matched, filter.Reversed), nil
	//round filter longitude in critical position(near to +-180)
	case !((limit/2 - math.Abs(float64(rfilter.CenterLatitude)) - float64(rfilter.Radius)) <= 0) && (limit/2-math.Abs(float64(rfilter.CentreLongitude))-float64(rfilter.Radius)) <= 0:
		matched := float64((pos.Latitude-rfilter.CenterLatitude)*(pos.Latitude-rfilter.CenterLatitude))+
			(limit-math.Abs(float64(pos.Longitude))-math.Abs(float64(rfilter.CentreLongitude)))*(limit-math.Abs(float64(pos.Longitude))-math.Abs(float64(rfilter.CentreLongitude))) <=
			float64(rfilter.Radius*rfilter.Radius)
		return xor(matched, filter.Reversed), nil
	//round filter canter in safe position (not near to 180)
	case !((180 - math.Abs(float64(rfilter.CenterLatitude)) - float64(rfilter.Radius)) <= 0) && !((180 - math.Abs(float64(rfilter.CentreLongitude)) - float64(rfilter.Radius)) <= 0):
		matched := (pos.Latitude-rfilter.CenterLatitude)*(pos.Latitude-rfilter.CenterLatitude)+
			(pos.Longitude-rfilter.CentreLongitude)*(pos.Longitude-rfilter.CentreLongitude) <=
			(rfilter.Radius * rfilter.Radius)
		return xor(matched, filter.Reversed), nil
	default:
		return false, nil
	}
}

func xor(a, b bool) bool {
	return (a && !b) || (!a && b)
}

// if leftLatitude > rightLatitude filter crosses twelve meridian and makes conflicts in matching
func isConflict(leftLatitude, rightLatitude float32) bool {
	return leftLatitude > rightLatitude
}

// moveRectangularFilter move filter right latitude to criticalRightLatitude
// and left latitude to point
// got by subtraction delta from original left latitude
func moveRectangularFilter(r *repository.RectangularFilter) float32 {
	delta := criticalRightLatitude + r.BottomRightLatitude
	r.BottomRightLatitude -= delta
	// move out from overlapping with degrees border from 180 to -180
	r.BottomRightLatitude += 360
	r.TopLeftLatitude -= delta

	return delta
}

// movePosition moves position latitude to point got by subtraction delta from position latitude
func movePosition(p *position.Position, delta float32) {
	p.Latitude -= delta
	if p.Latitude <= criticalLeftLatitude {
		// move out from overlapping with degrees border from 180 to -180
		p.Latitude += 360
	}
}

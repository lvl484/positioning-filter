// Package matcher provides matching positioning messages to a set of filtering rules
// and pereating messages to a topic in case of they are matched.
package matcher

import (
	"encoding/json"
	"errors"

	"github.com/lvl484/positioning-filter/position"
	"github.com/lvl484/positioning-filter/repository"
)

const (
	ErrBadFilterType = "Bad type of filter"
)

type matcher func(position.Position, *repository.Filter) (bool, error)

type matcherFilters struct {
	filters repository.Filters
}

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

	if !(rfilter.BottomRightLatitude > pos.Latitude) ||
		!(rfilter.TopLeftLatitude < pos.Latitude) ||
		!(rfilter.BottomRightLongitude < pos.Longitude) ||
		!(rfilter.TopLeftLongitude > pos.Longitude) {
		if filter.Reversed {
			return true, nil
		}

		return false, nil
	}

	if filter.Reversed {
		return false, nil
	}

	return true, nil
}

func matchRound(pos position.Position, filter *repository.Filter) (bool, error) {
	var rfilter repository.RoundFilter
	if err := json.Unmarshal(filter.Configuration, &rfilter); err != nil {
		return false, err
	}

	if (pos.Latitude-rfilter.CenterLatitude)*(pos.Latitude-rfilter.CenterLatitude)+
		(pos.Longitude-rfilter.CentreLongitude)*(pos.Longitude-rfilter.CentreLongitude) >
		(rfilter.Radius * rfilter.Radius) {
		if filter.Reversed {
			return true, nil
		}

		return false, nil
	}

	if filter.Reversed {
		return false, nil
	}

	return true, nil
}

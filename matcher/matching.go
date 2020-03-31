// Package matcher provides matching positioning messages to a set of filtering rules
// and pereating messages to a topic in case of they are matched.
package matcher

import (
	"github.com/lvl484/positioning-filter/position"
	"github.com/lvl484/positioning-filter/repository"
)

type Matcher interface {
	Match(position.Position) (bool, error)
}

type matcher struct {
	filters repository.Filters
}

func (m *matcher) Match(pos position.Position) (bool, error) {
	return true, nil
}

func NewMatcher(filters repository.Filters) Matcher {
	return &matcher{
		filters: filters,
	}
}

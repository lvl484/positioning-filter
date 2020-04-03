// Package matcher provides matching positioning messages to a set of filtering rules
// and pereating messages to a topic in case of they are matched.
package matcher

import (
	"github.com/lvl484/positioning-filter/position"
)

type Matcher interface {
	Match(position.Position) (bool, error)
}

package position

import (
	"time"

	"github.com/google/uuid"
)

type Position struct {
	UserID    uuid.UUID // id of user that produced position notification event
	Latitude  float32   // obviously latitude component of coordinates
	Longitude float32   // obviously longitude component of coordinates
	Timestamp time.Time // time when position collected
	Arrival   time.Time // time when position accepted by system
}

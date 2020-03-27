package position

import (
	"time"

	"github.com/google/uuid"
)

type Position struct {
	UserID    uuid.UUID `json:"UserID"`    // id of user that produced position notification event
	Latitude  float32   `json:"Latitude"`  // obviously latitude component of coordinates
	Longitude float32   `json:"Longitude"` // obviously longitude component of coordinates
	Timestamp time.Time `json:"Timestamp"` // time when position collected
	Arrival   time.Time `json:"Arrival"`   // time when position accepted by system
}

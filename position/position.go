package position

import (
	"time"

	"github.com/google/uuid"
)

type Position struct {
	UserID    uuid.UUID `json:"UserID"`
	Latitude  float32   `json:"Latitude"`
	Longitude float32   `json:"Longitude"`
	Timestamp time.Time `json:"Timestamp"`
	Arrival   time.Time `json:"Arrival"`
}

// Package repository provides filter model
package repository

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Filter is a struct to manipulate users' filter objects in database
type Filter struct {
	Name          string
	Type          string
	Configuration json.RawMessage
	Reversed      bool
	UserID        uuid.UUID
}

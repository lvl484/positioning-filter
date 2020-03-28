// Package repository provides filter model
package repository

// RoundFilter keeps configuration for filtering inside round area
type RoundFilter struct {
	CenterLatitude  float32
	CentreLongitude float32
	Radius          float32
}

// Package filter provides filter model
package filter

// RoundFilter keeps configuration for filtering inside round area
type RoundFilter struct {
	CenterLatitude  float32
	CentreLongitude float32
	Radius          float32
}

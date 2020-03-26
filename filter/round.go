// Package filter provides filter model
package filter

//RoundFilter detects if the position is inside of round area
type RoundFilter struct {
	CenterLatitude  float32
	CentreLongitude float32
	Radius          float32
}

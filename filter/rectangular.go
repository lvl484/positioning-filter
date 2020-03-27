// Package filter provides filter model
package filter

//RectangularFilter keeps configuration for filtering inside rectangular area
type RectangularFilter struct {
	TopLeftLatitude       float32
	TopLeftLongtitude     float32
	BottomRightLatitude   float32
	BottomRightLongtitude float32
}

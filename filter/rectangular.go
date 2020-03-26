// Package filter provides filter model
package filter

//RectangularFilter detects the position inside of rectangular area
type RectangularFilter struct {
	TopLeftLatitude       float32
	TopLeftLongtitude     float32
	BottomRightLatitude   float32
	BottomRightLongtitude float32
}

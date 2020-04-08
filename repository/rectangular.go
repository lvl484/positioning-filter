// Package repository provides filter model
package repository

//RectangularFilter keeps configuration for filtering inside rectangular area
type RectangularFilter struct {
	TopLeftLatitude      float32
	TopLeftLongitude     float32
	BottomRightLatitude  float32
	BottomRightLongitude float32
}

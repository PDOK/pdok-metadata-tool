package metadata

type BoundingBox struct {
	WestBoundLongitude string
	EastBoundLongitude string
	SouthBoundLatitude string
	NorthBoundLatitude string
}

type OnLine struct {
	URL      string
	Protocol string
}

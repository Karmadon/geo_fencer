package geo_fencer

import (
	"fmt"

	"github.com/karmadon/geo"

	"github.com/karmadon/geo_fencer/drivers"
)

type GeoFence interface {
	// Indexes this feature
	Add(f *geo.Feature)
	// Get all features that contain this coordinate
	Get(c geo.Coordinate) []*geo.Feature
}

func NewFence() GeoFence {
	return drivers.NewRtree()
}

func GetFence(label GeoFenceType, zoom int) (fence GeoFence, err error) {
	switch label {
	case GeoFenceTypeRtree:
		fence = drivers.NewRtree()
	case GeoFenceTypeS2:
		fence = drivers.NewS2fence(zoom)
	case GeoFenceTypeBbox:
		fence = drivers.NewBboxFence()
	case GeoFenceTypeBrute:
		fence = drivers.NewBruteFence()
	default:
		err = fmt.Errorf("bad fence type: %s", label)
	}
	return
}

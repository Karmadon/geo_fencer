package drivers

import (
	"github.com/stremovskyy/geo"
)

// Standard rtree with M=50
type Rtree struct {
	rtree *geo.Rtree
}

func NewRtree() *Rtree {
	return &Rtree{
		rtree: geo.NewRtree(),
	}
}

func (r *Rtree) Add(feature *geo.Feature) {
	for _, shape := range feature.Geometry {
		if len(shape.Coordinates) > 1 {
			r.rtree.Insert(geo.NewRnode(shape, feature))
		}
	}
}

func (r *Rtree) Get(coordinate geo.Coordinate) (matches []*geo.Feature) {
	nodes := r.rtree.Contains(coordinate)
	for _, node := range nodes {
		feature := node.Feature()
		if feature.Contains(coordinate) {
			matches = append(matches, feature)
		}
	}
	return
}

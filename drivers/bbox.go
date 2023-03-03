package drivers

import (
	"github.com/stremovskyy/geo"
)

type box struct {
	box     geo.Box
	feature *geo.Feature
}

type BboxFence struct {
	boxes []*box
}

func NewBboxFence() *BboxFence {
	return &BboxFence{}
}

func (b *BboxFence) Add(feature *geo.Feature) {
	for _, shape := range feature.Geometry {
		box := &box{box: shape.BoundingBox(), feature: feature}
		b.boxes = append(b.boxes, box)
	}
}

func (b *BboxFence) Get(c geo.Coordinate) (matches []*geo.Feature) {
	for _, box := range b.boxes {
		if box.box.Contains(c) {
			for _, shape := range box.feature.Geometry {
				if shape.Contains(c) {
					matches = append(matches, box.feature)
				}
			}
		}
	}
	return matches
}

package drivers

import "github.com/stremovskyy/geo"

type BruteFence struct {
	features []*geo.Feature
}

func NewBruteFence() *BruteFence {
	return &BruteFence{}
}

func (b *BruteFence) Add(feature *geo.Feature) {
	b.features = append(b.features, feature)
}

func (b *BruteFence) Get(c geo.Coordinate) (matches []*geo.Feature) {
	for _, feature := range b.features {
		for _, shape := range feature.Geometry {
			if shape.Contains(c) {
				matches = append(matches, feature)
			}
		}
	}
	return matches
}

func (b *BruteFence) Size() int {
	return len(b.features)
}

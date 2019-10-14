package drivers

import (
	"github.com/golang/geo/s2"
	"github.com/karmadon/geo"
)

type S2fence struct {
	resolution int
	covers     map[s2.CellID][]cover
}

func NewS2fence(resolution int) *S2fence {
	return &S2fence{
		resolution: resolution,
		covers:     make(map[s2.CellID][]cover),
	}
}

func (s *S2fence) Add(f *geo.Feature) {
	coverer := NewFlatCoverer(s.resolution)

	for _, shape := range f.Geometry {
		if shape.IsClockwise() {
			shape.Reverse()
		}

		points := make([]s2.Point, len(shape.Coordinates))
		for i, c := range shape.Coordinates[:] {
			points[i] = s2.PointFromLatLng(s2.LatLngFromDegrees(c.Lat, c.Lon))
		}

		region := s2.Region(LoopRegionFromPoints(points))

		bounds := coverer.Covering(region)
		if len(bounds) < 1 {
			continue
		}

		interiors := coverer.InteriorCovering(region)
		c := cover{
			interior: make(map[s2.CellID]bool, len(interiors)),
			feature:  f,
		}

		for _, cellID := range interiors {
			c.interior[cellID] = true
		}

		for _, cellID := range bounds {
			s.covers[cellID] = append(s.covers[cellID], c)
		}
	}
}

func (s *S2fence) Get(coordinate geo.Coordinate) (matches []*geo.Feature) {
	cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(coordinate.Lat, coordinate.Lon)).Parent(s.resolution)
	for _, cover := range s.covers[cellID] {
		if _, ok := cover.interior[cellID]; ok {
			matches = append(matches, cover.feature)
		} else if cover.feature.Contains(coordinate) {
			matches = append(matches, cover.feature)
		}
	}

	return
}

type cover struct {
	feature *geo.Feature
	interior map[s2.CellID]bool
}

type LoopRegion struct {
	*s2.Loop
}

func LoopRegionFromPoints(points []s2.Point) *LoopRegion {
	loop := s2.LoopFromPoints(points)
	return &LoopRegion{loop}
}

func (l *LoopRegion) CapBound() s2.Cap {
	return l.RectBound().CapBound()
}

func (l *LoopRegion) ContainsCell(cell s2.Cell) bool {
	for i := 0; i < 4; i++ {
		v := cell.Vertex(i)
		if !l.ContainsPoint(v) {
			return false
		}
	}
	return true
}

func (l *LoopRegion) IntersectsCell(cell s2.Cell) bool {
	for i := 0; i < 4; i++ {
		crosser := s2.NewChainEdgeCrosser(cell.Vertex(i), cell.Vertex((i+1)%4), l.Vertex(0))
		for _, verticle := range l.Vertices()[1:] {
			if crosser.EdgeOrVertexChainCrossing(verticle) {
				return true
			}
		}
		if crosser.EdgeOrVertexChainCrossing(l.Vertex(0)) { //close the loop
			return true
		}
	}
	return l.ContainsCell(cell)
}

type FlatCoverer struct {
	*s2.RegionCoverer
}

func NewFlatCoverer(level int) *FlatCoverer {
	return &FlatCoverer{&s2.RegionCoverer{
		MinLevel: level,
		MaxLevel: level,
		LevelMod: 0,
		MaxCells: 1 << 12,
	}}
}

func (c *FlatCoverer) Covering(r s2.Region) ( cover s2.CellUnion ){
	cellUnions := c.FastCovering(r.CapBound())
	for _, cellID := range cellUnions {
		cell := s2.CellFromCellID(cellID)
		if r.IntersectsCell(cell) {
			cover = append(cover, cellID)
		}
	}

	return cover
}

func (c *FlatCoverer) CellUnion(region s2.Region) s2.CellUnion {
	cover := c.Covering(region)
	cover.Normalize()
	return cover
}

func (c *FlatCoverer) InteriorCovering(region s2.Region) (cover s2.CellUnion) {
	cids := c.FastCovering(region.CapBound())
	for _, cid := range cids {
		cell := s2.CellFromCellID(cid)
		if region.ContainsCell(cell) {
			cover = append(cover, cid)
		}
	}
	return cover
}

func (c *FlatCoverer) InteriorCellUnion(region s2.Region) s2.CellUnion {
	cover := c.InteriorCovering(region)
	cover.Normalize()
	return cover
}

package geo_fencer

type GeoFenceType string

const (
	GeoFenceTypeS2    GeoFenceType = "s2"
	GeoFenceTypeBbox  GeoFenceType = "bbox"
	GeoFenceTypeBrute GeoFenceType = "brute"
	GeoFenceTypeRtree GeoFenceType = "rtree"
)

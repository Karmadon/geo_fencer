# geo_fencer - Simple Geo fence library

`geo_fencer` is a simple Go package for creating and searching geographic fences. It provides an interface for creating an index of multiple fences, making it easy to search for coordinates across different data sets.

## Installation

To install the package, use go get:
```shell
go get github.com/stremovskyy/geo_fencer
```

## Usage
To use geo_fencer, you need to first create a GeoFence object. You can then add features to the fence, which are essentially geographic shapes (such as polygons or circles) that define the boundaries of the fence.

Once you have created a fence, you can add it to a FenceIndex, which is a dictionary of multiple fences. You can then search for a coordinate in any of the indexed fences.

```go
import (
    "github.com/stremovskyy/geo_fencer"
    "github.com/stremovskyy/geo"
)

// Create a fence
fence, err := geo_fencer.GetFence(geo_fencer.GeoFenceTypeCircle, 10)
if err != nil {
    panic(err)
}

// Add a feature to the fence
center := geo.NewPoint(-122.431297, 37.773972)
radius := 1000.0
feature := geo.NewCircle(center, radius)
fence.Add(feature)

// Create an index of fences
fences := geo_fencer.NewFenceIndex()
fences.Set("my_fence", fence)

// Search for a coordinate
coord := geo.NewPoint(-122.431297, 37.773972)
match, err := fences.Search("my_fence", coord)
if err != nil {
    panic(err)
}
fmt.Println(match)

```

You can also load a fence index from a GeoJSON string using the LoadFenceIndex function:
```go
geoJsonString := `{
    "type": "FeatureCollection",
    "features": [
        {
            "type": "Feature",
            "geometry": {
                "type": "Polygon",
                "coordinates": [
                    [
                        [-122.4412, 37.7776],
                        [-122.4412, 37.7818],
                        [-122.4356, 37.7818],
                        [-122.4356, 37.7776],
                        [-122.4412, 37.7776]
                    ]
                ]
            }
        }
    ]
}`

fences, err := geo_fencer.LoadFenceIndex(geo_fencer.GeoFenceTypePolygon, 10, "my_fence", &geoJsonString)
if err != nil {
    panic(err)
}

coord := geo.NewPoint(-122.4384, 37.7797)
match, err := fences.Search("my_fence", coord)
if err != nil {
    panic(err)
}
fmt.Println(match)
```

## License

geo_fencer is licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.

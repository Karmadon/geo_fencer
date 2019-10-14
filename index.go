package geo_fencer

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/karmadon/geo"
	"github.com/paulmach/go.geojson"
)

//FenceIndex is a dictionary of multiple fences. Useful if you have multiple data sets that need to be searched
type FenceIndex interface {
	// Set the GeoFence
	Set(name string, fence GeoFence)
	// Get the GeoFence at the key, return nil if doesn't exist
	Get(name string) GeoFence
	// Add a feature to the GeoFence at the key
	Add(name string, feature *geo.Feature) error
	// Search for the coordinate at the key
	Search(name string, c geo.Coordinate) ([]*geo.Feature, error)
	// List the keys of the indexed fences
	Keys() []string
}

// Returns a thread-safe FenceIndex
func NewFenceIndex() FenceIndex {
	return NewMutexFenceIndex()
}

type UnsafeFenceIndex struct {
	fences map[string]GeoFence
}

func NewUnsafeFenceIndex() *UnsafeFenceIndex {
	return &UnsafeFenceIndex{fences: make(map[string]GeoFence)}
}

func (u *UnsafeFenceIndex) Set(name string, fence GeoFence) {
	u.fences[name] = fence
}

func (u *UnsafeFenceIndex) Get(name string) (fence GeoFence) {
	return u.fences[name]
}

func (u *UnsafeFenceIndex) Add(name string, feature *geo.Feature) (err error) {
	fence, ok := u.fences[name]
	if !ok {
		return fmt.Errorf("FenceIndex does not contain fence %q", name)
	}
	fence.Add(feature)
	return
}

func (u *UnsafeFenceIndex) Search(name string, c geo.Coordinate) (matchs []*geo.Feature, err error) {
	fence, ok := u.fences[name]
	if !ok {
		err = fmt.Errorf("FenceIndex does not contain fence %q", name)
		return
	}
	matchs = fence.Get(c)
	return
}

func (u *UnsafeFenceIndex) Keys() (keys []string) {
	for k := range u.fences {
		keys = append(keys, k)
	}
	return
}

type MutexFenceIndex struct {
	fences *UnsafeFenceIndex
	sync.RWMutex
}

func NewMutexFenceIndex() *MutexFenceIndex {
	return &MutexFenceIndex{fences: NewUnsafeFenceIndex()}
}

func (m *MutexFenceIndex) Set(name string, fence GeoFence) {
	m.Lock()
	defer m.Unlock()
	m.fences.Set(name, fence)
}

func (m *MutexFenceIndex) Get(name string) GeoFence {
	m.RLock()
	defer m.RUnlock()
	return m.fences.Get(name)
}

func (m *MutexFenceIndex) Add(name string, feature *geo.Feature) error {
	m.Lock()
	defer m.Unlock()
	return m.fences.Add(name, feature)
}

func (m *MutexFenceIndex) Search(name string, c geo.Coordinate) ([]*geo.Feature, error) {
	m.RLock()
	defer m.RUnlock()
	return m.fences.Search(name, c)
}

func (m *MutexFenceIndex) Keys() []string {
	m.RLock()
	defer m.RUnlock()
	return m.fences.Keys()
}

func LoadFenceIndex(fenceType GeoFenceType, resolution int, fenceName string, geoJsonString *string) (fences FenceIndex, err error) {
	if geoJsonString == nil {
		return nil, nil
	}

	fences = NewFenceIndex()
	fence, err := GetFence(fenceType, resolution)
	if err != nil {
		return nil, err
	}

	collection, err := geojson.UnmarshalFeatureCollection([]byte(*geoJsonString))
	if err != nil {
		return nil, err
	}

	for _, feature := range collection.Features {
		k, err := json.Marshal(feature)
		if err != nil {
			return nil, err
		}

		g, err := geo.UnmarshalGeojsonFeature(string(k))
		if err != nil {
			return nil, err
		}
		feature, err := geo.GeojsonFeatureAdapter(g)
		if err != nil {
			return nil, err
		}

		fence.Add(feature)
	}

	fences.Set(fenceName, fence)

	return fences, nil
}

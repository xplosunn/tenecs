package binding

import (
	"github.com/benbjohnson/immutable"
	"golang.org/x/exp/constraints"
)

// TwoLevelMap is like immutable.Map[S, immutable.Map[K, V]]
// but you can also add K, V pairs that are available for any S
type TwoLevelMap[S constraints.Ordered, K constraints.Ordered, V any] struct {
	scopedMaps *immutable.Map[S, *immutable.Map[K, V]]
	globalMap  *immutable.Map[K, V]
}

func NewTwoLevelMap[S constraints.Ordered, K constraints.Ordered, V any]() TwoLevelMap[S, K, V] {
	return TwoLevelMap[S, K, V]{
		scopedMaps: immutable.NewMap[S, *immutable.Map[K, V]](nil),
		globalMap:  immutable.NewMap[K, V](nil),
	}
}

func (m TwoLevelMap[S, K, V]) Get(scope S, key K) (V, bool) {
	value, ok := m.globalMap.Get(key)
	if ok {
		return value, ok
	}
	innerMap, ok := m.scopedMaps.Get(scope)
	if ok {
		return innerMap.Get(key)
	}
	var noResult V
	return noResult, false
}

func (m TwoLevelMap[S, K, V]) SetScopedIfAbsent(scope S, key K, value V) (TwoLevelMap[S, K, V], bool) {
	if _, ok := m.Get(scope, key); ok {
		return m, false
	}
	innerMap, ok := m.scopedMaps.Get(scope)
	if !ok {
		innerMap = immutable.NewMap[K, V](nil)
	}
	innerMap = innerMap.Set(key, value)
	return TwoLevelMap[S, K, V]{
		scopedMaps: m.scopedMaps.Set(scope, innerMap),
		globalMap:  m.globalMap.Delete(key),
	}, true
}

func (m TwoLevelMap[S, K, V]) SetGlobalIfAbsent(key K, value V) (TwoLevelMap[S, K, V], bool) {
	iterator := m.scopedMaps.Iterator()
	for !iterator.Done() {
		_, scopedMap, _ := iterator.Next()
		if _, ok := scopedMap.Get(key); ok {
			return m, false
		}
	}
	if _, ok := m.globalMap.Get(key); ok {
		return m, false
	}
	return TwoLevelMap[S, K, V]{
		scopedMaps: m.scopedMaps,
		globalMap:  m.globalMap.Set(key, value),
	}, true
}

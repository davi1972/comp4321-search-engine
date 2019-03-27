package concurrentMap

import(
	"sync"
)

// Map type that can be safely shared between
// goroutines that require read/write access to a map
// PageMap
type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// Concurrent map item
type ConcurrentMapItem struct {
	Key   string
	Value interface{}
}

func (cm *ConcurrentMap) Init() *ConcurrentMap {
	cm.items = make(map[string]interface{})
	return cm
}

// Sets a key in a concurrent map
func (cm *ConcurrentMap) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()

	cm.items[key] = value
}

// Gets a key from a concurrent map
func (cm *ConcurrentMap) Get(key string) (interface{}) {
	cm.Lock()
	defer cm.Unlock()

	value, _ := cm.items[key]

	return value
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *ConcurrentMap) Iter() <-chan ConcurrentMapItem {
	c := make(chan ConcurrentMapItem)

	f := func() {
		cm.Lock()
		defer cm.Unlock()

		for k, v := range cm.items {
			c <- ConcurrentMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}

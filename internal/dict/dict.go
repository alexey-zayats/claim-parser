package dict

import (
	"sync"
)

// Dict concurrent access map
type Dict struct {
	mutex *sync.RWMutex
	items map[string]interface{}
}

// New returns ready to use dict
func New() *Dict {
	return &Dict{
		mutex: &sync.RWMutex{},
		items: make(map[string]interface{}),
	}
}

// Get gets value from items by key
func (c *Dict) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

// Set value in items by key under Mutex
func (c *Dict) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = value
}

// Delete delete item by key under Mutex
func (c *Dict) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

// Foreach iterate over dict items and call callback function
func (c *Dict) Foreach(callback func(string, interface{})) {
	c.mutex.RLock()
	for key, elem := range c.items {
		callback(key, elem)
	}
	c.mutex.RUnlock()
}

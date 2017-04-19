//   Copyright 2015-2017 Ivan A Kostko (github.com/ivan-kostko; github.com/gopot)

//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at

//       http://www.apache.org/licenses/LICENSE-2.0

//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package concurrentmap

import "sync"

// The ConcurrentMap type represents light-weight and simple API for concurrent-safe operations over map[interface{}]interface{}
// Keys must be of comparable types as for general map[interface{}]interface{}
//
// NOTE(x): In case of operating on big amounts of data or need of extended functionality - consider to use https://github.com/streamrail/concurrent-map
type ConcurrentMap struct {
	items map[interface{}]interface{}
	lock  sync.RWMutex
}

// Private factory. It assigns items and set up RWMutex
func newConcurrentMap(items map[interface{}]interface{}) *ConcurrentMap {
	return &ConcurrentMap{items: items, lock: sync.RWMutex{}}
}

// Generic factory. Instantiates and initializes ConcurrentMap with `initCap` capacity.
func New(initCap int) *ConcurrentMap {
	items := make(map[interface{}]interface{}, initCap)
	return newConcurrentMap(items)
}

// Makes Concurrent copy of the `m`.
func MakeConcurrentCopy(m map[interface{}]interface{}) *ConcurrentMap {
	items := make(map[interface{}]interface{}, len(m))
	for key, value := range m {
		items[key] = value
	}
	return newConcurrentMap(items)
}

// Makes Concurrent copy of the `m` recursively.
// In case the value is map[string]interface{} - it converts it into ConcurrentMap recursively as well.
func MakeRecursivelyConcurrentCopy(m map[interface{}]interface{}) *ConcurrentMap {
	items := make(map[interface{}]interface{}, len(m))
	for key, value := range m {
		if x, ok := (value).(map[interface{}]interface{}); ok {
			items[key] = MakeRecursivelyConcurrentCopy(x)
		} else {
			items[key] = value
		}
	}
	return newConcurrentMap(items)
}

// Retrieves an element from map under given key.
// Returns false in case there is no entry associated with the key.
func (this *ConcurrentMap) Get(key interface{}) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	val, ok := this.items[key]
	return val, ok
}

// Sets the given value under the specified key.
func (this *ConcurrentMap) Set(key interface{}, val interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.items == nil {
		// we would need atleast one element in map
		this.items = make(map[interface{}]interface{}, 1)
	}

	this.items[key] = val
}

// Sets the given value under the specified key and returns true, if the key didn't exist uppon invokation.
// Returns false and does nothing, in case there is already an entry with the same key.
func (this *ConcurrentMap) SetIfNotExists(key interface{}, val interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.items[key]; !ok {
		if this.items == nil {
			// we would need atleast one element in map
			this.items = make(map[interface{}]interface{}, 1)
		}
		this.items[key] = val
		return true
	}
	return false
}

// Removes an element from the map.
func (this *ConcurrentMap) Remove(key interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.items, key)
}

// Returns copy of content as non concurrent(general) map.
func (this *ConcurrentMap) Items() map[interface{}]interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	x := make(map[interface{}]interface{}, len(this.items))
	for key, value := range this.items {
		x[key] = value
	}
	return x
}

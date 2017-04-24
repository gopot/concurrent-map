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

import "encoding/json"

// Implements [Unmarshaller](https://golang.org/pkg/encoding/json/#Unmarshaler).
//
// It recursively unmarshal values of JSON sructures into *ConcurrentMap, similar to unmarshalling into map[string]interface{}.
// Also, if some value represents a slice, it inspects its elements and unmarshals them into *ConcurrentMap if possible.
//
// While unmarshalling on non-empty map, overlapping key-values are overwritten.
//
// It is safe to concurrently Unmarshal, Get and/or Set. However, JSON key-value(s) are guaranteed to be available as up to date only by completion of UnmarshalJSON() method.
//
// TODO(gopot) : improve performance and minimize allocs.
func (this *ConcurrentMap) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	cm := convertIntoConcurrentMapRecursively(m)

	for key, value := range cm.Items() {
		this.Set(key, value)
	}

	return nil
}

func inspectAndConvertValueRecursively(in interface{}) interface{} {
	var out interface{}

	switch x := in.(type) {
	case map[string]interface{}:
		out = convertIntoConcurrentMapRecursively(x)
	case []map[string]interface{}:
		sl := make([]*ConcurrentMap, len(x))
		for i, j := range x {
			sl[i] = convertIntoConcurrentMapRecursively(j)
		}
		out = sl
	case []interface{}:
		sl := make([]interface{}, len(x))
		for i, j := range x {
			sl[i] = inspectAndConvertValueRecursively(j)
		}
		out = sl
	default:
		out = x
	}
	return out
}

func convertIntoConcurrentMapRecursively(m map[string]interface{}) *ConcurrentMap {
	cm := New(len(m))

	for key, value := range m {
		v := inspectAndConvertValueRecursively(value)

		cm.Set(key, v)
	}
	return cm
}

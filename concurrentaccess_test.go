//   Copyright © 2015-2017 Ivan A Kostko (github.com/ivan-kostko; github.com/gopot)

//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at

//       http://www.apache.org/licenses/LICENSE-2.0

//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// +build race

package concurrentmap_test

import (
	. "concurrent-map"
	"reflect"
	"sync"
	"testing"
)

func TestConcurrentOperations(t *testing.T) {
	testCases := []struct {
		TestAlias   string
		Iterations  int
		Routines    int
		Cm          *ConcurrentMap
		Func        func(cm *ConcurrentMap, N int)
		ExpectedMap *ConcurrentMap
	}{
		{
			TestAlias:  `SetIfNotExists,Get,Set,Remove,Items OneKey`,
			Iterations: 1000,
			Routines:   16,
			Cm:         new(ConcurrentMap),
			Func: func(cm *ConcurrentMap, N int) {
				for i := 0; i < N; i++ {
					// initial kv
					const key, value = 0, 0

					cm.SetIfNotExists(key, value)

					if iv, ok := cm.Get(key); ok {
						if v, ok := iv.(int); ok {
							cm.Set(key, v+1)
						} else {
							cm.Set(key, value)
						}
					}

					if len(cm.Items()) > 0 {
						cm.Remove(key)
					}
				}
			},
			ExpectedMap: New(0),
		},
		{
			TestAlias:  `SetIfNotExists,Get,Set,Remove,Items Sequential Keys`,
			Iterations: 1000,
			Routines:   128,
			Cm:         new(ConcurrentMap),
			Func: func(cm *ConcurrentMap, N int) {
				for i := 0; i < N; i++ {

					cm.SetIfNotExists(i, i)

					if v, ok := cm.Get(i); ok {
						if iv, ok := v.(int); ok {
							cm.Set(i, N-iv)
						}
					}

					// just to invoke cm.Items()
					if len(cm.Items()) > 0 {
						cm.Remove(i)
					}
				}
			},
			ExpectedMap: New(0),
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		iterations := testCase.Iterations
		routines := testCase.Routines
		cm := testCase.Cm
		fn := testCase.Func
		expectedMap := testCase.ExpectedMap

		testFn := func(t *testing.T) {

			wg := sync.WaitGroup{}

			for r := 0; r < routines; r++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fn(cm, iterations)
				}()
			}

			wg.Wait()

			if !(reflect.DeepEqual(cm, expectedMap)) {
				t.Errorf("%s :: \r\n returned map is \r\n %#v \r\n while expected %#v", testAlias, cm, expectedMap)
			}
		}

		t.Run(testAlias, testFn)

	}

}

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

package concurrentmap_test

import (
	. "concurrent-map"
	"reflect"
	"testing"
)

func TestNewItemsCycle(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		InitCapacity  int
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "Converting simple map into a concurrent map",
			InitCapacity:  1,
			ExpectedItems: map[interface{}]interface{}{},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		initCapacity := testCase.InitCapacity
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			newCm := New(initCapacity)

			actualItems := newCm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: New(%d) TSM with items as \r\n %v \r\n while expected \r\n %v ", testAlias, initCapacity, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestMakeConcurrentCopyItemsCycle(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		OriginalMap   map[interface{}]interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "Converting simple map into a concurrent map",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			copyCm := MakeConcurrentCopy(originalMap)

			actualItems := copyCm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: MakeConcurrentCopy TSM with items as \r\n %v \r\n while expected \r\n %v ", testAlias, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestMakeRecursivelyConcurrentCopyItemsCycle(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		OriginalMap   map[interface{}]interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "Converting simple map into a concurrent map",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
		{
			TestAlias:     "Converting nested map into a concurrent map",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "Map": map[interface{}]interface{}{"key1": "stringValue", "key2": 123}},
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "Map": MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123})},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			copyCm := MakeRecursivelyConcurrentCopy(originalMap)

			actualItems := copyCm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: (MakeRecursivelyConcurrentCopy(%#v)).Items() returned \r\n %#v \r\n while expected \r\n %#v ", testAlias, originalMap, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestGet(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		Cm            *ConcurrentMap
		GetAKey       interface{}
		ExpectedValue interface{}
		ExpectedOk    bool
	}{
		{
			TestAlias:     "MakeConcurrentCopy and Get existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			GetAKey:       "key2",
			ExpectedValue: 123,
			ExpectedOk:    true,
		},
		{
			TestAlias:     "MakeConcurrentCopy and Get non-existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			GetAKey:       "key3",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Get existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			GetAKey:       "key2",
			ExpectedValue: 123,
			ExpectedOk:    true,
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Get non-existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			GetAKey:       "key3",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			TestAlias:     "New(0) and Get non-existing key",
			Cm:            New(0),
			GetAKey:       "key3",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			TestAlias:     "new(ConcurrentMap) and Get non-existing key",
			Cm:            new(ConcurrentMap),
			GetAKey:       "key3",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		getAKey := testCase.GetAKey
		expectedValue := testCase.ExpectedValue
		expectedOk := testCase.ExpectedOk

		testFn := func(t *testing.T) {

			actualValue, actualOk := cm.Get(getAKey)

			if !(reflect.DeepEqual(actualValue, expectedValue)) {
				t.Errorf("%s :: cm.Get(%s) returned Value \r\n %#v \r\n while expected \r\n %#v ", testAlias, getAKey, actualValue, expectedValue)
			}
			if actualOk != expectedOk {
				t.Errorf("%s :: cm.Get(%s) returned error \r\n %#v \r\n while expected \r\n %#v ", testAlias, getAKey, actualOk, expectedOk)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestSetGetCycle(t *testing.T) {

	testCases := []struct {
		TestAlias string
		Cm        *ConcurrentMap
		SetAKey   interface{}
		SetAValue interface{}
	}{
		{
			TestAlias: "MakeConcurrentCopy and Set existing key",
			Cm:        MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:   "key2",
			SetAValue: 321,
		},
		{
			TestAlias: "MakeConcurrentCopy and Set non-existing key",
			Cm:        MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:   "key3",
			SetAValue: 4.56,
		},
		{
			TestAlias: "MakeRecursivelyConcurrentCopy and Set existing key",
			Cm:        MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:   "key2",
			SetAValue: 321,
		},
		{
			TestAlias: "MakeRecursivelyConcurrentCopy and Set non-existing key",
			Cm:        MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:   "key3",
			SetAValue: 4.56,
		},
		{
			TestAlias: "New(0) and Set non-existing key",
			Cm:        New(0),
			SetAKey:   "key3",
			SetAValue: 4.56,
		},
		{
			TestAlias: "new(ConcurrentMap) and Set non-existing key",
			Cm:        new(ConcurrentMap),
			SetAKey:   "key3",
			SetAValue: 4.56,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue

		testFn := func(t *testing.T) {

			cm.Set(setAKey, setAValue)

			actualValue, actualOk := cm.Get(setAKey)

			if !(reflect.DeepEqual(actualValue, setAValue)) {
				t.Errorf("%s :: cm.Get('%s') after cm.Set('%s', %#v) returned Value \r\n %#v \r\n while expected \r\n %#v ", testAlias, setAKey, setAKey, setAValue, actualValue, setAValue)
			}
			if !actualOk {
				t.Errorf("%s :: cm.Get('%s') after cm.Set('%s', %#v) returned OK as \r\n %#v \r\n while expected true ", testAlias, setAKey, setAKey, setAValue, actualOk)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestSetItemsCycle(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		Cm            *ConcurrentMap
		SetAKey       interface{}
		SetAValue     interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "MakeConcurrentCopy and Set existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key2",
			SetAValue:     321,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 321},
		},
		{
			TestAlias:     "MakeConcurrentCopy and Set non-existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "key3": 4.56},
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Set existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key2",
			SetAValue:     321,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 321},
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Set non-existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "key3": 4.56},
		},
		{
			TestAlias:     "New(0) and Set non-existing key",
			Cm:            New(0),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key3": 4.56},
		},
		{
			TestAlias:     "new(ConcurrentMap) and Set non-existing key",
			Cm:            new(ConcurrentMap),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key3": 4.56},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			cm.Set(setAKey, setAValue)

			actualItems := cm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: cm.Items() after cm.Set('%s', %#v) returned \r\n %#v \r\n while expected \r\n %#v ", testAlias, setAKey, setAValue, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestSetIfNotExistsReturnOk(t *testing.T) {

	testCases := []struct {
		TestAlias  string
		Cm         *ConcurrentMap
		SetAKey    interface{}
		SetAValue  interface{}
		ExpectedOk bool
	}{
		{
			TestAlias:  "MakeConcurrentCopy and Set existing key",
			Cm:         MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:    "key2",
			SetAValue:  321,
			ExpectedOk: false,
		},
		{
			TestAlias:  "MakeConcurrentCopy and Set non-existing key",
			Cm:         MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:    "key3",
			SetAValue:  4.56,
			ExpectedOk: true,
		},
		{
			TestAlias:  "MakeRecursivelyConcurrentCopy and Set existing key",
			Cm:         MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:    "key2",
			SetAValue:  321,
			ExpectedOk: false,
		},
		{
			TestAlias:  "MakeRecursivelyConcurrentCopy and Set non-existing key",
			Cm:         MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:    "key3",
			SetAValue:  4.56,
			ExpectedOk: true,
		},
		{
			TestAlias:  "New(0) and Set non-existing key",
			Cm:         New(0),
			SetAKey:    "key3",
			SetAValue:  4.56,
			ExpectedOk: true,
		},
		{
			TestAlias:  "new(ConcurrentMap) and Set non-existing key",
			Cm:         new(ConcurrentMap),
			SetAKey:    "key3",
			SetAValue:  4.56,
			ExpectedOk: true,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedOk := testCase.ExpectedOk

		testFn := func(t *testing.T) {

			actualOk := cm.SetIfNotExists(setAKey, setAValue)

			if actualOk != expectedOk {
				t.Errorf("%s :: cm.SetIfNotExists('%s', %#v) returned OK \r\n %#v \r\n while expected \r\n %#v ", testAlias, setAKey, setAValue, actualOk, expectedOk)
			}

		}
		t.Run(testAlias, testFn)
	}

}

func TestSetIfNotExistsGetCycle(t *testing.T) {

	testCases := []struct {
		TestAlias        string
		Cm               *ConcurrentMap
		SetAKey          interface{}
		SetAValue        interface{}
		ExpectedGetValue interface{}
	}{
		{
			TestAlias:        "MakeConcurrentCopy and Set existing key",
			Cm:               MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:          "key2",
			SetAValue:        321,
			ExpectedGetValue: 123,
		},
		{
			TestAlias:        "MakeConcurrentCopy and Set non-existing key",
			Cm:               MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:          "key3",
			SetAValue:        4.56,
			ExpectedGetValue: 4.56,
		},
		{
			TestAlias:        "MakeRecursivelyConcurrentCopy and Set existing key",
			Cm:               MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:          "key2",
			SetAValue:        321,
			ExpectedGetValue: 123,
		},
		{
			TestAlias:        "MakeRecursivelyConcurrentCopy and Set non-existing key",
			Cm:               MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:          "key3",
			SetAValue:        4.56,
			ExpectedGetValue: 4.56,
		},
		{
			TestAlias:        "New(0) and Set non-existing key",
			Cm:               New(0),
			SetAKey:          "key3",
			SetAValue:        4.56,
			ExpectedGetValue: 4.56,
		},
		{
			TestAlias:        "new(ConcurrentMap) and Set non-existing key",
			Cm:               new(ConcurrentMap),
			SetAKey:          "key3",
			SetAValue:        4.56,
			ExpectedGetValue: 4.56,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedGetValue := testCase.ExpectedGetValue

		testFn := func(t *testing.T) {

			_ = cm.SetIfNotExists(setAKey, setAValue)

			actualGetValue, actualGetOk := cm.Get(setAKey)

			if !(reflect.DeepEqual(actualGetValue, expectedGetValue)) {
				t.Errorf("%s :: cm.Get('%s') after cm.Set('%s', %#v) returned Value \r\n %#v \r\n while expected \r\n %#v ", testAlias, setAKey, setAKey, setAValue, actualGetValue, expectedGetValue)
			}
			if !actualGetOk {
				t.Errorf("%s :: cm.Get('%s') after cm.Set('%s', %#v) returned OK as \r\n %#v \r\n while expected true ", testAlias, setAKey, setAKey, setAValue, actualGetOk)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func TestSetIfNotExistsItemsCycle(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		Cm            *ConcurrentMap
		SetAKey       interface{}
		SetAValue     interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "MakeConcurrentCopy and Set existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key2",
			SetAValue:     321,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
		{
			TestAlias:     "MakeConcurrentCopy and Set non-existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "key3": 4.56},
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Set existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key2",
			SetAValue:     321,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Set non-existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "key3": 4.56},
		},
		{
			TestAlias:     "New(0) and Set non-existing key",
			Cm:            New(0),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key3": 4.56},
		},
		{
			TestAlias:     "new(ConcurrentMap) and Set non-existing key",
			Cm:            new(ConcurrentMap),
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key3": 4.56},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			_ = cm.SetIfNotExists(setAKey, setAValue)

			actualItems := cm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: cm.Items() after cm.SetIfNotExists('%s', %#v) returned \r\n %#v \r\n while expected \r\n %#v ", testAlias, setAKey, setAValue, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestRemoveGetCycle(t *testing.T) {

	testCases := []struct {
		TestAlias        string
		Cm               *ConcurrentMap
		RemoveAKey       interface{}
		ExpectedGetValue interface{}
		ExpectedGetOk    bool
	}{
		{
			TestAlias:        "MakeConcurrentCopy and Remove existing key",
			Cm:               MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:       "key2",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
		{
			TestAlias:        "MakeConcurrentCopy and Remove non-existing key",
			Cm:               MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:       "key3",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
		{
			TestAlias:        "MakeRecursivelyConcurrentCopy and Remove existing key",
			Cm:               MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:       "key2",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
		{
			TestAlias:        "MakeRecursivelyConcurrentCopy and Remove non-existing key",
			Cm:               MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:       "key3",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
		{
			TestAlias:        "New(0) and Remove non-existing key",
			Cm:               New(0),
			RemoveAKey:       "key3",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
		{
			TestAlias:        "new(ConcurrentMap) and Remove non-existing key",
			Cm:               new(ConcurrentMap),
			RemoveAKey:       "key3",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		removeAKey := testCase.RemoveAKey
		expectedGetValue := testCase.ExpectedGetValue
		expectedGetOk := testCase.ExpectedGetOk

		testFn := func(t *testing.T) {

			cm.Remove(removeAKey)

			actualGetValue, actualGetOk := cm.Get(removeAKey)

			if !(reflect.DeepEqual(actualGetValue, expectedGetValue)) {
				t.Errorf("%s :: cm.Get('%s') after cm.Remove('%s') returned Value \r\n %#v \r\n while expected \r\n %#v ", testAlias, removeAKey, removeAKey, actualGetValue, expectedGetValue)
			}
			if actualGetOk != expectedGetOk {
				t.Errorf("%s :: cm.Get('%s') after cm.Remove('%s') returned OK as \r\n %#v \r\n while expected %#v ", testAlias, removeAKey, removeAKey, actualGetOk, expectedGetOk)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func TestRemoveItemsCycle(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		Cm            *ConcurrentMap
		RemoveAKey    interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "MakeConcurrentCopy and Remove existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:    "key2",
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue"},
		},
		{
			TestAlias:     "MakeConcurrentCopy and Remove non-existing key",
			Cm:            MakeConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:    "key3",
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Remove existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:    "key2",
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue"},
		},
		{
			TestAlias:     "MakeRecursivelyConcurrentCopy and Remove non-existing key",
			Cm:            MakeRecursivelyConcurrentCopy(map[interface{}]interface{}{"key1": "stringValue", "key2": 123}),
			RemoveAKey:    "key3",
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
		{
			TestAlias:     "New(0) and Remove non-existing key",
			Cm:            New(0),
			RemoveAKey:    "key3",
			ExpectedItems: map[interface{}]interface{}{},
		},
		{
			TestAlias:     "new(ConcurrentMap) and Remove non-existing key",
			Cm:            new(ConcurrentMap),
			RemoveAKey:    "key3",
			ExpectedItems: map[interface{}]interface{}{},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		cm := testCase.Cm
		removeAKey := testCase.RemoveAKey
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			cm.Remove(removeAKey)

			actualItems := cm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: cm.Items() after cm.Remove('%s') returned \r\n %#v \r\n while expected \r\n %#v ", testAlias, removeAKey, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}
}

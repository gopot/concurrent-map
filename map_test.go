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

package concurrentmap_test

import (
	. "github.com/gopot/concurrent-map"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	testCases := []struct {
		TestAlias    string
		InitCapacity int
		ExpectedMap  map[interface{}]interface{}
	}{
		{
			TestAlias:    "Converting simple map into a concurrent map",
			InitCapacity: 1,
			ExpectedMap:  map[interface{}]interface{}{},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		initCapacity := testCase.InitCapacity
		expectedMap := testCase.ExpectedMap

		testFn := func(t *testing.T) {

			actualCm := New(initCapacity)

			if !(reflect.DeepEqual(actualCm.Items(), expectedMap)) {
				t.Errorf("%s :: New(%d) TSM with items as \r\n %v \r\n while expected \r\n %v ", testAlias, initCapacity, actualCm.Items(), expectedMap)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestMakeConcurrentCopy(t *testing.T) {

	testCases := []struct {
		TestAlias   string
		OriginalMap map[interface{}]interface{}
		ExpectedMap map[interface{}]interface{}
	}{
		{
			TestAlias:   "Converting simple map into a concurrent map",
			OriginalMap: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			ExpectedMap: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		expectedMap := testCase.ExpectedMap

		testFn := func(t *testing.T) {

			actualCm := MakeConcurrentCopy(originalMap)

			if !(reflect.DeepEqual(actualCm.Items(), expectedMap)) {
				t.Errorf("%s :: MakeConcurrentCopy TSM with items as \r\n %v \r\n while expected \r\n %v ", testAlias, actualCm.Items(), expectedMap)
			}
		}
		t.Run(testAlias, testFn)
	}

}

func TestMakeRecursivelyConcurrentCopy(t *testing.T) {

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

			actualCm := MakeRecursivelyConcurrentCopy(originalMap)

			actualItems := actualCm.Items()

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
		OriginalMap   map[interface{}]interface{}
		GetAKey       interface{}
		ExpectedValue interface{}
		ExpectedOk    bool
	}{
		{
			TestAlias:     "Get existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			GetAKey:       "key2",
			ExpectedValue: 123,
			ExpectedOk:    true,
		},
		{
			TestAlias:     "Get non-existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			GetAKey:       "key3",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		getAKey := testCase.GetAKey
		expectedValue := testCase.ExpectedValue
		expectedOk := testCase.ExpectedOk

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		TestAlias   string
		OriginalMap map[interface{}]interface{}
		SetAKey     interface{}
		SetAValue   interface{}
	}{
		{
			TestAlias:   "Set existing key",
			OriginalMap: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:     "key2",
			SetAValue:   321,
		},
		{
			TestAlias:   "Set non-existing key",
			OriginalMap: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:     "key3",
			SetAValue:   4.56,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		OriginalMap   map[interface{}]interface{}
		SetAKey       interface{}
		SetAValue     interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "Set existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:       "key2",
			SetAValue:     321,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 321},
		},
		{
			TestAlias:     "Set non-existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "key3": 4.56},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		TestAlias   string
		OriginalMap map[interface{}]interface{}
		SetAKey     interface{}
		SetAValue   interface{}
		ExpectedOk  bool
	}{
		{
			TestAlias:   "Set existing key",
			OriginalMap: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:     "key2",
			SetAValue:   321,
			ExpectedOk:  false,
		},
		{
			TestAlias:   "Set non-existing key",
			OriginalMap: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:     "key3",
			SetAValue:   4.56,
			ExpectedOk:  true,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedOk := testCase.ExpectedOk

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		OriginalMap      map[interface{}]interface{}
		SetAKey          interface{}
		SetAValue        interface{}
		ExpectedGetValue interface{}
	}{
		{
			TestAlias:        "Set existing key",
			OriginalMap:      map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:          "key2",
			SetAValue:        321,
			ExpectedGetValue: 123,
		},
		{
			TestAlias:        "Set non-existing key",
			OriginalMap:      map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:          "key3",
			SetAValue:        4.56,
			ExpectedGetValue: 4.56,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedGetValue := testCase.ExpectedGetValue

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		OriginalMap   map[interface{}]interface{}
		SetAKey       interface{}
		SetAValue     interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "Set existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:       "key2",
			SetAValue:     321,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
		{
			TestAlias:     "Set non-existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			SetAKey:       "key3",
			SetAValue:     4.56,
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123, "key3": 4.56},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		setAKey := testCase.SetAKey
		setAValue := testCase.SetAValue
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		OriginalMap      map[interface{}]interface{}
		RemoveAKey       interface{}
		ExpectedGetValue interface{}
		ExpectedGetOk    bool
	}{
		{
			TestAlias:        "Remove existing key",
			OriginalMap:      map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			RemoveAKey:       "key2",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
		{
			TestAlias:        "Remove non-existing key",
			OriginalMap:      map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			RemoveAKey:       "key3",
			ExpectedGetValue: nil,
			ExpectedGetOk:    false,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		removeAKey := testCase.RemoveAKey
		expectedGetValue := testCase.ExpectedGetValue
		expectedGetOk := testCase.ExpectedGetOk

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

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
		OriginalMap   map[interface{}]interface{}
		RemoveAKey    interface{}
		ExpectedItems map[interface{}]interface{}
	}{
		{
			TestAlias:     "Remove existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			RemoveAKey:    "key2",
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue"},
		},
		{
			TestAlias:     "Remove non-existing key",
			OriginalMap:   map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
			RemoveAKey:    "key3",
			ExpectedItems: map[interface{}]interface{}{"key1": "stringValue", "key2": 123},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalMap := testCase.OriginalMap
		removeAKey := testCase.RemoveAKey
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {

			cm := MakeConcurrentCopy(originalMap)

			cm.Remove(removeAKey)

			actualItems := cm.Items()

			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("%s :: cm.Items() after cm.Remove('%s') returned \r\n %#v \r\n while expected \r\n %#v ", testAlias, removeAKey, actualItems, expectedItems)
			}
		}
		t.Run(testAlias, testFn)
	}
}

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
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		InitialMap    *ConcurrentMap
		JsonData      []byte
		ExpectedError error
		ExpectedItems map[interface{}]interface{}
	}{

		{
			TestAlias:     "Simple key-value on empty map",
			InitialMap:    New(1),
			JsonData:      []byte(`{"key": "value"}`),
			ExpectedError: nil,
			ExpectedItems: map[interface{}]interface{}{"key": "value"},
		},
		{
			TestAlias:     "Nested key-value on empty map",
			InitialMap:    New(1),
			JsonData:      []byte(`{"key": {"key": "value"}}`),
			ExpectedError: nil,
			ExpectedItems: map[interface{}]interface{}{"key": MakeConcurrentCopy(map[interface{}]interface{}{"key": "value"})},
		},
		{
			TestAlias:     "Nested slice key-value on empty map",
			InitialMap:    New(1),
			JsonData:      []byte(`{"key": [{"key": "value"}, {"key": "value"}, {"key": "value"}]}`),
			ExpectedError: nil,
			ExpectedItems: map[interface{}]interface{}{
				"key": []interface{}{
					MakeConcurrentCopy(map[interface{}]interface{}{"key": "value"}),
					MakeConcurrentCopy(map[interface{}]interface{}{"key": "value"}),
					MakeConcurrentCopy(map[interface{}]interface{}{"key": "value"}),
				},
			},
		},
		{
			TestAlias:     "Complex nested slice key-value on empty map",
			InitialMap:    New(1),
			JsonData:      []byte(`{"key": [{"key1": "value"}, [{"key2": "value"}, {"key2": "value"}, {"key2": "value"}], {"key3": "value"}]}`),
			ExpectedError: nil,
			ExpectedItems: map[interface{}]interface{}{
				"key": []interface{}{
					MakeConcurrentCopy(map[interface{}]interface{}{"key1": "value"}),
					[]interface{}{
						MakeConcurrentCopy(map[interface{}]interface{}{"key2": "value"}),
						MakeConcurrentCopy(map[interface{}]interface{}{"key2": "value"}),
						MakeConcurrentCopy(map[interface{}]interface{}{"key2": "value"}),
					},
					MakeConcurrentCopy(map[interface{}]interface{}{"key3": "value"}),
				},
			},
		},
		{
			TestAlias:     "Complex nested slice of key-value slices on empty map",
			InitialMap:    New(1),
			JsonData:      []byte(`{"key": [[{"key1": "value"},[{"key21": "value"}, {"key22": "value"}, {"key23": "value"}]],[{"key11": "value"},[{"key12": "value"}, {"key13": "value"}, {"key11": "value"}],[{"key31": "value"},[{"key32": "value"}, {"key33": "value"}, {"key34": "value"}]], {"key3": "value"}]]}`),
			ExpectedError: nil,
			ExpectedItems: map[interface{}]interface{}{
				"key": []interface{}{
					[]interface{}{
						MakeConcurrentCopy(map[interface{}]interface{}{"key1": "value"}),

						[]interface{}{
							MakeConcurrentCopy(map[interface{}]interface{}{"key21": "value"}),
							MakeConcurrentCopy(map[interface{}]interface{}{"key22": "value"}),
							MakeConcurrentCopy(map[interface{}]interface{}{"key23": "value"}),
						},
					},
					[]interface{}{
						MakeConcurrentCopy(map[interface{}]interface{}{"key11": "value"}),
						[]interface{}{
							MakeConcurrentCopy(map[interface{}]interface{}{"key12": "value"}),
							MakeConcurrentCopy(map[interface{}]interface{}{"key13": "value"}),
							MakeConcurrentCopy(map[interface{}]interface{}{"key11": "value"}),
						},
						[]interface{}{
							MakeConcurrentCopy(map[interface{}]interface{}{"key31": "value"}),
							[]interface{}{
								MakeConcurrentCopy(map[interface{}]interface{}{"key32": "value"}),
								MakeConcurrentCopy(map[interface{}]interface{}{"key33": "value"}),
								MakeConcurrentCopy(map[interface{}]interface{}{"key34": "value"}),
							},
						},
						MakeConcurrentCopy(map[interface{}]interface{}{"key3": "value"}),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		initialMap := testCase.InitialMap
		jsonData := testCase.JsonData
		expectedError := testCase.ExpectedError
		expectedItems := testCase.ExpectedItems

		testFn := func(t *testing.T) {
			actualError := initialMap.UnmarshalJSON(jsonData)

			actualItems := initialMap.Items()

			if !(reflect.DeepEqual(actualError, expectedError)) {
				t.Errorf("initialMap.UnmarshalJSON(%s) \r\n returned error \r\n %+v \r\n while expected \r\n %+v \r\n", jsonData, actualError, expectedError)
			}
			if !(reflect.DeepEqual(actualItems, expectedItems)) {
				t.Errorf("initialMap.UnmarshalJSON(%s); initialMap.Items() \r\n returned \r\n %#v \r\n while expected \r\n %#v \r\n", jsonData, actualItems, expectedItems)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Benchmark_UnmarshalJSON_ConcurrentMap_vs_MapStringInterface(b *testing.B) {

	benchCases := []struct {
		TestAlias string
		JsonData  []byte
	}{

		{
			TestAlias: "Simple key-value on empty map",
			JsonData:  []byte(`{"key": "value"}`),
		},
		{
			TestAlias: "Nested key-value on empty map",
			JsonData:  []byte(`{"key": {"key": "value"}}`),
		},
		{
			TestAlias: "Nested slice key-value on empty map",
			JsonData:  []byte(`{"key": [{"key": "value"}, {"key": "value"}, {"key": "value"}]}`),
		},
		{
			TestAlias: "Complex nested slice key-value on empty map",
			JsonData:  []byte(`{"key": [{"key1": "value"}, [{"key2": "value"}, {"key2": "value"}, {"key2": "value"}], {"key3": "value"}]}`),
		},
		{
			TestAlias: "Complex nested slice of key-value slices on empty map",
			JsonData:  []byte(`{"key": [[{"key1": "value"},[{"key21": "value"}, {"key22": "value"}, {"key23": "value"}]],[{"key11": "value"},[{"key12": "value"}, {"key13": "value"}, {"key11": "value"}],[{"key31": "value"},[{"key32": "value"}, {"key33": "value"}, {"key34": "value"}]], {"key3": "value"}]]}`),
		},
	}

	for _, benchCase := range benchCases {
		testAlias := benchCase.TestAlias
		jsonData := benchCase.JsonData

		benchCmFn := func(b *testing.B) {
			b.ReportAllocs()

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				cm := New(0)
				b.StartTimer()
				json.Unmarshal(jsonData, cm)
				b.StopTimer()
			}
		}

		b.Run(`CM `+testAlias, benchCmFn)

		benchMapFn := func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				m := make(map[string]interface{})
				b.StartTimer()
				json.Unmarshal(jsonData, &m)
				b.StopTimer()
			}
		}

		b.Run(`Map `+testAlias, benchMapFn)

	}

}

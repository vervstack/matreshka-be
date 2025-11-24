package environment

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/evon"
	"gopkg.in/yaml.v3"
)

func Test_IntVariable(t *testing.T) {
	t.Parallel()

	const (
		singleIntValue   = 1
		singleInt8Value  = int8(singleIntValue)
		singleInt16Value = int16(singleIntValue)
		singleInt32Value = int32(singleIntValue)
		singleInt64Value = int64(singleIntValue)
	)

	type testCase struct {
		val      any
		expected int
	}

	testCases := map[string]testCase{
		"int": {
			val:      singleIntValue,
			expected: singleIntValue,
		},
		"int8": {
			val:      singleInt8Value,
			expected: singleIntValue,
		},
		"int16": {
			val:      singleInt16Value,
			expected: singleIntValue,
		},
		"int32": {
			val:      singleInt32Value,
			expected: singleIntValue,
		},
		"int64": {
			val:      singleInt64Value,
			expected: singleIntValue,
		},

		"negative_int": {
			val:      -singleIntValue,
			expected: -singleIntValue,
		},
		"negative_int8": {
			val:      -singleInt8Value,
			expected: -singleIntValue,
		},
		"negative_int16": {
			val:      -singleInt16Value,
			expected: -singleIntValue,
		},
		"negative_int32": {
			val:      -singleInt32Value,
			expected: -singleIntValue,
		},
		"negative_int64": {
			val:      -singleInt64Value,
			expected: -singleIntValue,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := MustNewVariable(varName, tc.val)

			expected := &Variable{
				Name: varName,
				Type: VariableTypeInt,
				Value: Value{
					val: &intValue{
						v: tc.expected,
					},
				},
			}

			require.Equal(t, expected, actual)

			marshalled, err := yaml.Marshal(actual)
			require.NoError(t, err)
			var expectedYaml string
			if tc.expected > 0 {
				expectedYaml = `
name: test_variable
type: int
value: 1
`
			} else {
				expectedYaml = `
name: test_variable
type: int
value: -1
`
			}

			require.YAMLEq(t, string(marshalled), expectedYaml)

			var newVar Variable
			err = yaml.Unmarshal(marshalled, &newVar)
			require.NoError(t, err)

			require.Equal(t, expected, &newVar)
		})
	}

}

func TestIntSliceVariable(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val  any
		opts []opt

		expected []int
	}

	var (
		intSlice         = []int{1, 3, 5}
		negativeIntSlice = []int{-1, -3, 5}
	)

	testCases := map[string]testCase{
		"int": {
			val:      []int{1, 3, 5},
			expected: intSlice,
		},
		"int8": {
			val:      []int8{1, 3, 5},
			expected: intSlice,
		},
		"int16": {
			val:      []int16{1, 3, 5},
			expected: intSlice,
		},
		"int32": {
			val:      []int32{1, 3, 5},
			expected: intSlice,
		},
		"int64": {
			val:      []int64{1, 3, 5},
			expected: intSlice,
		},
		"negative_int": {
			val:      []int{-1, -3, 5},
			expected: negativeIntSlice,
		},
		"negative_int8": {
			val:      []int8{-1, -3, 5},
			expected: negativeIntSlice,
		},
		"negative_int16": {
			val:      []int16{-1, -3, 5},
			expected: negativeIntSlice,
		},
		"negative_int32": {
			val:      []int32{-1, -3, 5},
			expected: negativeIntSlice,
		},
		"negative_int64": {
			val:      []int64{-1, -3, 5},
			expected: negativeIntSlice,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := MustNewVariable(varName, tc.val, tc.opts...)

			expected := &Variable{
				Name: varName,
				Type: VariableTypeInt,
				Value: Value{
					val: &intSliceValue{
						v: tc.expected,
					},
				},
			}

			require.Equal(t, expected, actual)

			marshalled, err := yaml.Marshal(actual)
			require.NoError(t, err)

			var expectedYaml string
			if tc.expected[0] > 0 {
				expectedYaml = `
name: test_variable
type: int
value: [1, 3, 5]
`
			} else {
				expectedYaml = `
name: test_variable
type: int
value: [-1, -3, 5]
`
			}

			require.YAMLEq(t, expectedYaml, string(marshalled))

			var newVar Variable
			err = yaml.Unmarshal(marshalled, &newVar)
			require.NoError(t, err)

			require.Equal(t, expected, &newVar)
		})
	}
}

func TestIntRangeVariable(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val  any
		opts []opt

		expected []int
	}

	positiveExpected := []int{1, 2, 3, 5, 6, 7}
	negativeExpected := []int{-7, -6, -5, -3, -2, -1}
	mixedExpected := []int{-2, -1, 0, 1, 3, 4, 5}

	testCases := map[string]testCase{
		"int": {
			val:      []int{1, 2, 3, 5, 6, 7},
			expected: positiveExpected,
		},
		"int8": {
			val:      []int8{1, 2, 3, 5, 6, 7},
			expected: positiveExpected,
		},
		"int16": {
			val:      []int16{1, 2, 3, 5, 6, 7},
			expected: positiveExpected,
		},
		"int32": {
			val:      []int32{1, 2, 3, 5, 6, 7},
			expected: positiveExpected,
		},
		"int64": {
			val:      []int64{1, 2, 3, 5, 6, 7},
			expected: positiveExpected,
		},

		"negative_int": {
			val:      []int{-7, -6, -5, -3, -2, -1},
			expected: negativeExpected,
		},
		"negative_int8": {
			val:      []int8{-7, -6, -5, -3, -2, -1},
			expected: negativeExpected,
		},
		"negative_int16": {
			val:      []int16{-7, -6, -5, -3, -2, -1},
			expected: negativeExpected,
		},
		"negative_int32": {
			val:      []int32{-7, -6, -5, -3, -2, -1},
			expected: negativeExpected,
		},
		"negative_int64": {
			val:      []int64{-7, -6, -5, -3, -2, -1},
			expected: negativeExpected,
		},

		"mixed_int": {
			val:      []int{-2, -1, 0, 1, 3, 4, 5},
			expected: mixedExpected,
		},
		"mixed_int8": {
			val:      []int8{-2, -1, 0, 1, 3, 4, 5},
			expected: mixedExpected,
		},

		"mixed_int16": {
			val:      []int16{-2, -1, 0, 1, 3, 4, 5},
			expected: mixedExpected,
		},
		"mixed_int32": {
			val:      []int32{-2, -1, 0, 1, 3, 4, 5},
			expected: mixedExpected,
		},
		"mixed_int64": {
			val:      []int64{-2, -1, 0, 1, 3, 4, 5},
			expected: mixedExpected,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := MustNewVariable(varName, tc.val, tc.opts...)

			expected := &Variable{
				Name: varName,
				Type: VariableTypeInt,
				Value: Value{
					val: &intSliceValue{
						v: tc.expected,
					},
				},
			}

			require.Equal(t, expected, actual)

			marshalled, err := yaml.Marshal(actual)
			require.NoError(t, err)
			var expectedYaml string
			if tc.expected[2] > 0 {
				expectedYaml = `
name: test_variable
type: int
value: [1:3, 5:7]
`
			} else if tc.expected[2] < 0 {
				expectedYaml = `
name: test_variable
type: int
value: [-7:-5, -3:-1]
`
			} else {
				expectedYaml = `
name: test_variable
type: int
value: [-2:1, 3:5]
`
			}

			require.YAMLEq(t, expectedYaml, string(marshalled))

			actualUnmarshalled := &Variable{}
			err = yaml.Unmarshal(marshalled, actualUnmarshalled)
			require.NoError(t, err)

			require.Equal(t, expected, actualUnmarshalled)

			//	Evon
			marshalledEvon, err := evon.MarshalEnv(actual)
			require.NoError(t, err)
			_ = marshalledEvon
		})
	}
}

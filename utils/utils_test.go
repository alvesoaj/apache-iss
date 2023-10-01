package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CastStringToUint8(t *testing.T) {
	testCases := map[string]struct {
		sVal     string
		base     int
		expected uint8
	}{
		"base 10": {
			sVal:     "100",
			base:     10,
			expected: 100,
		},
		"base 2": {
			sVal:     "1100100",
			base:     2,
			expected: 100,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			val := CastStringToUint8(testCase.sVal, testCase.base)
			assert.Equal(t, testCase.expected, val)
		})
	}
}

func Test_CastInterfaceToUint8(t *testing.T) {
	var i interface{} = uint8(100)
	val := CastInterfaceToUint8(i)
	assert.Equal(t, uint8(100), val)
}

func Test_RemoveAllNonNumericFromString(t *testing.T) {
	val := RemoveAllNonNumericFromString("> 100\n")
	assert.Equal(t, "100", val)
}

func Test_NewTestInput(t *testing.T) {
	in, err := NewTestInput("")
	assert.NoError(t, err)
	assert.NotNil(t, in)
}

func Test_NewTestOutput(t *testing.T) {
	out := NewTestOutput()
	assert.NotNil(t, out)
}

func Test_ClearOutputForTesting(t *testing.T) {
	val := ClearOutputForTesting("> 100")
	assert.Equal(t, "100", val)
}

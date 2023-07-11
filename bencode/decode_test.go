package bencode

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		input     string
		expected  map[string]interface{}
		expectErr bool
	}{
		{
			input:    "d3:key5:valuee",
			expected: map[string]interface{}{"key": "value"},
		},
		{
			input:    "d4:listli1ei2ee5:helloi42ee",
			expected: map[string]interface{}{"list": []interface{}{int64(1), int64(2)}, "hello": int64(42)},
		},
		{
			input:     "d3:keyi1234e",
			expectErr: true, // Invalid value type for key
		},
	}

	for _, tc := range testCases {
		reader := strings.NewReader(tc.input)
		decoder := NewDecoder(reader)

		result, err := decoder.Decode()
		if tc.expectErr && err == nil {
			t.Errorf("Expected an error, but got nil")
		} else if !tc.expectErr && err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, but got %v", tc.expected, result)
		}
	}
}

func TestDecodeError(t *testing.T) {
	invalidInput := "d3:key5:value" // Missing 'e' at the end
	reader := strings.NewReader(invalidInput)
	decoder := NewDecoder(reader)

	_, err := decoder.Decode()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

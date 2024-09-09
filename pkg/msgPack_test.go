package msgPack

import (
	"reflect"
	"strings"
	"testing"
)

func TestSerializer(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  []byte
		wantError bool
	}{
		{
			name:      "int8",
			input:     int8(127),
			expected:  []byte("d07f"),
			wantError: false,
		},
		{
			name:      "int16",
			input:     int16(1024),
			expected:  []byte("d10400"),
			wantError: false,
		},
		{
			name:      "int32",
			input:     int32(65536),
			expected:  []byte("d200010000"),
			wantError: false,
		},
		{
			name:      "int64",
			input:     int64(1234567890),
			expected:  []byte("d300000000499602d2"),
			wantError: false,
		},
		{
			name:      "uint8",
			input:     uint8(200),
			expected:  []byte("ccc8"),
			wantError: false,
		},
		{
			name:      "uint16",
			input:     uint16(1024),
			expected:  []byte("cd0400"),
			wantError: false,
		},
		{
			name:      "uint32",
			input:     uint32(65536),
			expected:  []byte("ce00010000"),
			wantError: false,
		},
		{
			name:      "uint64",
			input:     uint64(1234567890),
			expected:  []byte("cf00000000499602d2"),
			wantError: false,
		},
		{
			name:      "nil",
			input:     nil,
			expected:  []byte("c0"),
			wantError: false,
		},
		{
			name:      "bool true",
			input:     true,
			expected:  []byte("c3"),
			wantError: false,
		},
		{
			name:      "bool false",
			input:     false,
			expected:  []byte("c2"),
			wantError: false,
		},
		{
			name:      "float32",
			input:     float32(3.14),
			expected:  []byte("ca4048f5c3"),
			wantError: false,
		},
		{
			name:      "float64",
			input:     float64(3.14),
			expected:  []byte("cb40091eb851eb851f"),
			wantError: false,
		},
		{
			name:      "string less than 16 characters",
			input:     "hello",
			expected:  []byte("a5hello"),
			wantError: false,
		},
		{
			name:      "string between 16 and 31 characters",
			input:     "this is a longer string",
			expected:  []byte("b7this is a longer string"),
			wantError: false,
		},
		{
			name:      "string more than 31 characters",
			input:     "this string has more than 31 characters",
			expected:  []byte("d927this string has more than 31 characters"),
			wantError: false,
		},
		{
			name:      "string more than 2^16 -1 characters",
			input:     strings.Repeat("a", (1<<16)-1),
			expected:  []byte("daffff" + strings.Repeat("a", (1<<16)-1)),
			wantError: false,
		},
		{
			name:      "array of different types",
			input:     []interface{}{int8(1), uint8(255), "test"},
			expected:  []byte("93d001ccffa4test"),
			wantError: false,
		},
		{
			name:      "unsupported data type",
			input:     struct{}{},
			wantError: true,
		},
		{
			name:      "string more than 2^32 characters",
			input:     string(make([]byte, 1<<32)),
			wantError: true,
		},
		{
			name:      "empty array",
			input:     []interface{}{},
			expected:  []byte("90"),
			wantError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Serializer(test.input)
			if (err != nil) != test.wantError {
				t.Errorf("error got = %e, wantError= %v", err, test.wantError)
				return
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got = %v, expected: %v", got, test.expected)
			}
		})
	}
}

func TestDeserializer(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expected  interface{}
		wantError bool
	}{
		{
			name:      "int8",
			input:     []byte("d07f"),
			expected:  int64(127),
			wantError: false,
		},
		{
			name:      "int16",
			input:     []byte("d10400"),
			expected:  int64(1024),
			wantError: false,
		},
		{
			name:      "int32",
			input:     []byte("d200010000"),
			expected:  int64(65536),
			wantError: false,
		},
		{
			name:      "int64",
			input:     []byte("d300000000499602d2"),
			expected:  int64(1234567890),
			wantError: false,
		},
		{
			name:      "uint8",
			input:     []byte("ccc8"),
			expected:  uint64(200),
			wantError: false,
		},
		{
			name:      "uint16",
			input:     []byte("cd0400"),
			expected:  uint64(1024),
			wantError: false,
		},
		{
			name:      "uint32",
			input:     []byte("ce00010000"),
			expected:  uint64(65536),
			wantError: false,
		},
		{
			name:      "uint64",
			input:     []byte("cf00000000499602d2"),
			expected:  uint64(1234567890),
			wantError: false,
		},
		{
			name:      "nil",
			input:     []byte("c0"),
			expected:  nil,
			wantError: false,
		},
		{
			name:      "bool true",
			input:     []byte("c3"),
			expected:  true,
			wantError: false,
		},
		{
			name:      "bool false",
			input:     []byte("c2"),
			expected:  false,
			wantError: false,
		},
		{
			name:      "float32",
			input:     []byte("ca4048f5c3"),
			expected:  float64(3.140000104904175),
			wantError: false,
		},
		{
			name:      "float64",
			input:     []byte("cb40091eb851eb851f"),
			expected:  float64(3.14),
			wantError: false,
		},
		{
			name:      "string less than 16 characters",
			input:     []byte("a5hello"),
			expected:  "hello",
			wantError: false,
		},
		{
			name:      "string between 16 and 31 characters",
			input:     []byte("b7this is a longer string"),
			expected:  "this is a longer string",
			wantError: false,
		},
		{
			name:      "string more than 31 characters",
			input:     []byte("d927this string has more than 31 characters"),
			expected:  "this string has more than 31 characters",
			wantError: false,
		},
		{
			name:      "string more than 2^16 -1 characters",
			input:     append([]byte("daffff"), []byte(strings.Repeat("a", (1<<16)-1))...),
			expected:  strings.Repeat("a", (1<<16)-1),
			wantError: false,
		},
		{
			name:      "array of different types",
			input:     []byte("93d001ccffa4test"),
			expected:  []interface{}{int64(1), uint64(255), "test"},
			wantError: false,
		},
		{
			name:      "unsupported data type",
			input:     []byte("ww"),
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Deserializer(&test.input)
			if (err != nil) != test.wantError {
				t.Errorf("error got = %e, wantError= %v", err, test.wantError)
				return
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got = %v, expected = %v", got, test.expected)
			}
		})
	}
}

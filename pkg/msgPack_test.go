package msgPack

import (
	"reflect"
	"strings"
	"testing"
)

func TestSerializer(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
		wantErr  bool
	}{
		{
			name:     "int8",
			input:    int8(127),
			expected: []byte("d07f"),
			wantErr:  false,
		},
		{
			name:     "int16",
			input:    int16(1024),
			expected: []byte("d10400"),
			wantErr:  false,
		},
		{
			name:     "int32",
			input:    int32(65536),
			expected: []byte("d200010000"),
			wantErr:  false,
		},
		{
			name:     "int64",
			input:    int64(1234567890),
			expected: []byte("d300000000499602d2"),
			wantErr:  false,
		},
		{
			name:     "uint8",
			input:    uint8(200),
			expected: []byte("ccc8"),
			wantErr:  false,
		},
		{
			name:     "uint16",
			input:    uint16(1024),
			expected: []byte("cd0400"),
			wantErr:  false,
		},
		{
			name:     "uint32",
			input:    uint32(65536),
			expected: []byte("ce00010000"),
			wantErr:  false,
		},
		{
			name:     "uint64",
			input:    uint64(1234567890),
			expected: []byte("cf00000000499602d2"),
			wantErr:  false,
		},
		{
			name:     "nil",
			input:    nil,
			expected: []byte("c0"),
			wantErr:  false,
		},
		{
			name:     "bool true",
			input:    true,
			expected: []byte("c3"),
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: []byte("c2"),
			wantErr:  false,
		},
		{
			name:     "float32",
			input:    float32(3.14),
			expected: []byte("ca4048f5c3"),
			wantErr:  false,
		},
		{
			name:     "float64",
			input:    float64(3.14),
			expected: []byte("cb40091eb851eb851f"),
			wantErr:  false,
		},
		{
			name:     "string less than 16 characters",
			input:    "hello",
			expected: []byte("a5hello"),
			wantErr:  false,
		},
		{
			name:     "string between 16 and 31 characters",
			input:    "this is a longer string",
			expected: []byte("b7this is a longer string"),
			wantErr:  false,
		},
		{
			name:     "string more than 31 characters",
			input:    "this string has more than 31 characters",
			expected: []byte("d927this string has more than 31 characters"),
			wantErr:  false,
		},
		{
			name:     "string more than 2^16 -1 characters",
			input:    strings.Repeat("a",(1<<16) -1),
			expected: []byte("daffff"+strings.Repeat("a",(1<<16) -1)),
			wantErr:  false,
		},
		{
			name:     "array of different types",
			input:    []interface{}{int8(1), uint8(255), "test"},
			expected: []byte("93d001ccffa4test"),
			wantErr:  false,
		},
		{
			name:    "unsupported data type",
			input:   struct{}{},
			wantErr: true,
		},
		{
			name:    "string more than 2^32 characters",
			input:   string(make([]byte, 1<<32)),
			wantErr: true,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: []byte("90"),
			wantErr:  false,
		},
		// {
		// 	name:     "map of string to interface",
		// 	input:    map[string]interface{}{"key1": "value1", "key2": "value2"},
		// 	expected: []byte("82a4key2a6value2a4key1a6value1"),
		// 	wantErr:  false,
		// },
		// {
		// 	name:     "byte slice",
		// 	input:    []byte{0x01, 0x02, 0x03},
		// 	expected: []byte("c403010203"),
		// 	wantErr:  false,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Serializer(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("Serializer() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("Serializer() = %x, want %x", got, test.expected)
			}
		})
	}
}


func TestDeserializer(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "int8",
			input:    []byte("d07f"),
			expected: int64(127),
			wantErr:  false,
		},
		{
			name:     "int16",
			input:    []byte("d10400"),
			expected: int64(1024),
			wantErr:  false,
		},
		{
			name:     "int32",
			input:    []byte("d200010000"),
			expected: int64(65536),
			wantErr:  false,
		},
		{
			name:     "int64",
			input:    []byte("d300000000499602d2"),
			expected: int64(1234567890),
			wantErr:  false,
		},
		{
			name:     "uint8",
			input:    []byte("ccc8"),
			expected: uint64(200),
			wantErr:  false,
		},
		{
			name:     "uint16",
			input:    []byte("cd0400"),
			expected: uint64(1024),
			wantErr:  false,
		},
		{
			name:     "uint32",
			input:    []byte("ce00010000"),
			expected: uint64(65536),
			wantErr:  false,
		},
		{
			name:     "uint64",
			input:    []byte("cf00000000499602d2"),
			expected: uint64(1234567890),
			wantErr:  false,
		},
		{
			name:     "nil",
			input:    []byte("c0"),
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "bool true",
			input:    []byte("c3"),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    []byte("c2"),
			expected: false,
			wantErr:  false,
		},
		{
			name:     "float32",
			input:    []byte("ca4048f5c3"),
			expected: float64(3.140000104904175),
			wantErr:  false,
		},
		{
			name:     "float64",
			input:    []byte("cb40091eb851eb851f"),
			expected: float64(3.14),
			wantErr:  false,
		},
		{
			name:     "string less than 16 characters",
			input:    []byte("a5hello"),
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "string between 16 and 31 characters",
			input:    []byte("b7this is a longer string"),
			expected: "this is a longer string",
			wantErr:  false,
		},
		{
			name:     "string more than 31 characters",
			input:    []byte("d927this string has more than 31 characters"),
			expected: "this string has more than 31 characters",
			wantErr:  false,
		},
		{
			name:     "string more than 2^16 -1 characters",
			input:    append([]byte("daffff"), []byte(strings.Repeat("a", (1<<16)-1))...),
			expected: strings.Repeat("a", (1<<16)-1),
			wantErr:  false,
		},
		{
			name:     "array of different types",
			input:    []byte("93d001ccffa4test"),
			expected: []interface{}{int64(1), uint64(255), "test"},
			wantErr:  false,
		},
		{
			name:    "unsupported data type",
			input:   []byte("ww"),
			wantErr: true,
		},
		// {
		// 	name:     "map of string to interface",
		// 	input:    []byte("82a4key2a6value2a4key1a6value1"),
		// 	expected: map[string]interface{}{"key1": "value1", "key2": "value2"},
		// 	wantErr:  false,
		// },
		// {
		// 	name:     "byte slice",
		// 	input:    []byte("c403010203"),
		// 	expected: []byte{0x01, 0x02, 0x03},
		// 	wantErr:  false,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Deserializer(&test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("Deserializer() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("Deserializer() = %v, want %v", got, test.expected)
			}
		})
	}
}
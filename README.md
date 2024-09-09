# Msgpack-Nabil

Encoder and Decode of msgPAck

## Installation

As a library

```shell
go get github.com/codescalersinternships/Msgpack-Nabil/pkg
```

## Usage

in your Go app you can do something like

```go
package main

import (
	"fmt"

	msgPack "github.com/codescalersinternships/Msgpack-Nabil/pkg"
)


func main() {
	encoded, err := msgPack.Serializer([]interface{}{
		int8(127), int16(32767), int32(2147483647), int64(9223372036854775807),
		uint8(255), uint16(65535), uint32(4294967295), uint64(18446744073709551615),
		float32(3.14159), float64(2.718281828459045),
		"short string", "this is a longer string that is used to test the string handling capability of the serializer",
		[]byte{0x00, 0xFF, 0x7F},
		true, false,
		map[string]interface{}{"key": "value"},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(encoded))
	fmt.Println(msgPack.Deserializer(&encoded))
}

```

## Testing

```shell
make test
```
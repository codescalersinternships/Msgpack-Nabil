package msgPack

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func numByted(num int64, cnt int, encoded *[]byte) {
	numStr := strconv.FormatInt(num, 16)
	for len(numStr) < cnt/4 {
		numStr = "0" + numStr
	}
	for st, en := 0, 2; st < cnt/4; st += 2 {
		*encoded = append(*encoded, []byte(numStr[st:en])...)
		en += 2
	}
}

func numBytedUnsigned(num uint64, cnt int, encoded *[]byte) {
	numStr := strconv.FormatUint(num, 16)
	for len(numStr) < cnt/4 {
		numStr = "0" + numStr
	}
	for st, en := 0, 2; st < cnt/4; st += 2 {
		*encoded = append(*encoded, []byte(numStr[st:en])...)
		en += 2
	}
}

// Serializer interface to MessagePack bytes and return array of bytes
// if there is error it will return error
func Serializer(applicationObject interface{}) ([]byte, error) {
	applicationObjectCopy := applicationObject
	return serializer(applicationObjectCopy)
}

func serializer(applicationObject interface{}) ([]byte, error) {
	var encoded []byte

	switch ty := applicationObject.(type) {
	case int8:
		num := reflect.ValueOf(ty).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD0, 16))...)
		numByted(num, 8, &encoded)

	case int16:
		num := reflect.ValueOf(ty).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD1, 16))...)
		numByted(num, 16, &encoded)

	case int32:
		num := reflect.ValueOf(ty).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD2, 16))...)
		numByted(num, 32, &encoded)

	case int64:
		num := reflect.ValueOf(ty).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD3, 16))...)
		numByted(num, 64, &encoded)

	case uint8:
		num := reflect.ValueOf(ty).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCC, 16))...)
		numBytedUnsigned(num, 8, &encoded)

	case uint16:
		num := reflect.ValueOf(ty).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCD, 16))...)
		numBytedUnsigned(num, 16, &encoded)

	case uint32:
		num := reflect.ValueOf(ty).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCE, 16))...)
		numBytedUnsigned(num, 32, &encoded)

	case uint64:
		num := reflect.ValueOf(ty).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCF, 16))...)
		numBytedUnsigned(num, 64, &encoded)

	case nil:
		encoded = append(encoded, []byte(strconv.FormatUint(0xC0, 16))...)

	case bool:
		val := reflect.ValueOf(ty).Bool()
		if val {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC3, 16))...)
		} else {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC2, 16))...)
		}

	case float32:
		encoded = append(encoded, []byte(strconv.FormatUint(0xCA, 16))...)
		bits := math.Float32bits(ty)
		enc := make([]byte, 4)
		binary.BigEndian.PutUint32(enc, bits)
		hexString := hex.EncodeToString(enc)
		encoded = append(encoded, []byte(hexString)...)

	case float64:
		encoded = append(encoded, []byte(strconv.FormatUint(0xCB, 16))...)
		bits := math.Float64bits(ty)
		enc := make([]byte, 8)
		binary.BigEndian.PutUint64(enc, bits)
		hexString := hex.EncodeToString(enc)
		encoded = append(encoded, []byte(hexString)...)

	case string:
		val := reflect.ValueOf(ty).String()
		if len(val) < 32 {
			if len(val) < 16 {
				encoded = append(encoded, []byte("a"+strconv.FormatInt(int64(len(val)), 16))...)
			} else {
				encoded = append(encoded, []byte("b"+strconv.FormatInt(int64(len(val))-16, 16))...)
			}
		} else if len(val) < (1 << 8) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xD9, 16))...)
			numBytedUnsigned(uint64(len(val)), 8, &encoded)
		} else if len(val) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDA, 16))...)
			numBytedUnsigned(uint64(len(val)), 16, &encoded)
		} else if len(val) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDB, 16))...)
			numBytedUnsigned(uint64(len(val)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		encoded = append(encoded, []byte(val)...)

	case []byte:
		val := reflect.ValueOf(ty).Bytes()
		if len(val) < (1 << 8) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC4, 16))...)
			numBytedUnsigned(uint64(len(val)), 8, &encoded)
		} else if len(val) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC5, 16))...)
			numBytedUnsigned(uint64(len(val)), 16, &encoded)
		} else if len(val) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC6, 16))...)
			numBytedUnsigned(uint64(len(val)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		encoded = append(encoded, []byte(val)...)

	case []interface{}:
		if len(ty) < (1 << 4) {
			numBytedUnsigned(uint64(0x90+len(ty)), 8, &encoded)
		} else if len(ty) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDC, 16))...)
			numBytedUnsigned(uint64(len(ty)), 16, &encoded)
		} else if len(ty) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDD, 16))...)
			numBytedUnsigned(uint64(len(ty)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		for _, val := range ty {
			elements, err := serializer(val)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, elements...)
		}

	case map[string]interface{}:
		if len(ty) < (1 << 4) {
			numBytedUnsigned(uint64(0x80+len(ty)), 8, &encoded)
		} else if len(ty) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDE, 16))...)
			numBytedUnsigned(uint64(len(ty)), 16, &encoded)
		} else if len(ty) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDF, 16))...)
			numBytedUnsigned(uint64(len(ty)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		for key, val := range ty {
			elementKey, err := serializer(key)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, elementKey...)
			elementVal, err := serializer(val)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, elementVal...)
		}
	default:
		return nil, fmt.Errorf("unsupported data type")
	}

	return encoded, nil

}

// Deserialze MessagePack bytes and return interface
// if there is error it will return error
func Deserializer(encoded *[]byte) (interface{}, error) {
	encodedCopy := encoded
	return deserializer(encodedCopy)
}

func deserializer(encoded *[]byte) (interface{}, error) {

	var decoded interface{}

	var err error
	x, err := parseUintN((*encoded)[0:2], 8)
	if err != nil {
		return nil, fmt.Errorf("unsupported data type")
	}
	*encoded = (*encoded)[2:]
	switch x {
	case 0xD0: // int8
		decoded, err = parseIntN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[2:]
	case 0xD1: // int16
		decoded, err = parseIntN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[4:]
	case 0xD2: // int32
		decoded, err = parseIntN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
	case 0xD3: // int64
		decoded, err = parseIntN((*encoded)[:16], 64)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[16:]
	case 0xCC: // uint8
		decoded, err = parseUintN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[2:]
	case 0xCD: // uint16
		decoded, err = parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[4:]
	case 0xCE: // uint32
		decoded, err = parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
	case 0xCF: // uint64
		decoded, err = parseUintN((*encoded)[:16], 64)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[16:]
	case 0xC0: // nil
		decoded = nil
	case 0xC2: // false
		decoded = false
	case 0xC3: // true
		decoded = true
	case 0xCA: // float32
		decoded, err = parseFloatN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
	case 0xCB: // float64
		decoded, err = parseFloatN((*encoded)[:16], 64)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[16:]
	case 0xD9: // string 8
		length, err := parseUintN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		decoded = string((*encoded)[2 : 2+length])
		*encoded = (*encoded)[2+length:]
	case 0xDA: // string 16
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		decoded = string((*encoded)[4 : 4+length])
		*encoded = (*encoded)[4+length:]
	case 0xDB: // string 32
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		decoded = string((*encoded)[8 : 8+length])
		*encoded = (*encoded)[8+length:]
	case 0xC4: // bin 8
		length, err := parseUintN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		decoded = (*encoded)[2 : 2+length]
		*encoded = (*encoded)[2+length:]
	case 0xC5: // bin 16
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		decoded = (*encoded)[4 : 4+length]
		*encoded = (*encoded)[4+length:]
	case 0xC6: // bin 32
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		decoded = (*encoded)[8 : 8+length]
		*encoded = (*encoded)[8+length:]
	case 0xDC: // array 16
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[4:]
		decoded, err = parseArray(encoded, int(length))
		if err != nil {
			return nil, err
		}
	case 0xDD: // array 32
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
		decoded, err = parseArray(encoded, int(length))
		if err != nil {
			return nil, err
		}
	case 0xDE: // map 16
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[4:]
		decoded, err = parseMap(encoded, int(length))
		if err != nil {
			return nil, err
		}
	case 0xDF: // map 32
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
		decoded, err = parseMap(encoded, int(length))
		if err != nil {
			return nil, err
		}
	default:
		typ := x
		if typ <= 0x7F {
			decoded = uint8(typ)
		} else if typ <= 0x8f {
			decoded, err = parseMap(encoded, int(typ)-int(0x80))
			if err != nil {
				return nil, err
			}
		} else if typ <= 0x9f {
			decoded, err = parseArray(encoded, int(typ)-int(0x90))
			if err != nil {
				return nil, err
			}
		} else if typ <= 0xbf {
			length := typ - 0xa0
			decoded = string((*encoded)[:length])
			*encoded = (*encoded)[length:]
		} else if typ >= 0xe0 && typ <= 0xff {
			decoded = int8(typ)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
	}

	return decoded, nil
}

func parseIntN(encoded []byte, bitSize int) (int64, error) {
	return strconv.ParseInt(string(encoded), 16, bitSize)
}

func parseUintN(encoded []byte, bitSize int) (uint64, error) {
	return strconv.ParseUint(string(encoded), 16, bitSize)
}

func parseFloatN(encoded []byte, bitSize int) (float64, error) {
	hexStr := string(encoded)
	decodedBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, err
	}
	switch bitSize {
	case 32:
		bits := binary.BigEndian.Uint32(decodedBytes)
		return float64(math.Float32frombits(bits)), nil
	case 64:
		bits := binary.BigEndian.Uint64(decodedBytes)
		return math.Float64frombits(bits), nil
	default:
		return 0, fmt.Errorf("unsupported float bit size: %d", bitSize)
	}
}

func parseArray(encoded *[]byte, length int) ([]interface{}, error) {
	var arr []interface{}
	var err error
	for i := 0; i < length; i++ {
		var elem interface{}
		elem, err = deserializer(encoded)
		if err != nil {
			return nil, err
		}
		arr = append(arr, elem)
	}
	return arr, nil
}

func parseMap(encoded *[]byte, length int) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for i := 0; i < length; i++ {
		key, err := deserializer(encoded)
		if err != nil {
			return nil, err
		}
		value, err := deserializer(encoded)
		if err != nil {
			return nil, err
		}
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("key is not a string")
		}
		m[keyStr] = value
	}
	return m, nil
}

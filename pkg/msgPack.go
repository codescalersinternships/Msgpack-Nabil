package msgPack

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func numByted(num int64, nibbleCnt int, encodedRef *[]byte) {
	numStr := strconv.FormatInt(num, 16)
	for len(numStr) < nibbleCnt/4 {
		numStr = "0" + numStr
	}
	for st, en := 0, 2; st < nibbleCnt/4; st += 2 {
		*encodedRef = append(*encodedRef, []byte(numStr[st:en])...)
		en += 2
	}
}

func numBytedUnsigned(num uint64, nibbleCnt int, encodedRef *[]byte) {
	numStr := strconv.FormatUint(num, 16)
	for len(numStr) < nibbleCnt/4 {
		numStr = "0" + numStr
	}
	for st, en := 0, 2; st < nibbleCnt/4; st += 2 {
		*encodedRef = append(*encodedRef, []byte(numStr[st:en])...)
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

	switch curAppObject := applicationObject.(type) {
	case int8:
		num := reflect.ValueOf(curAppObject).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD0, 16))...)
		numByted(num, 8, &encoded)

	case int16:
		num := reflect.ValueOf(curAppObject).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD1, 16))...)
		numByted(num, 16, &encoded)

	case int32:
		num := reflect.ValueOf(curAppObject).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD2, 16))...)
		numByted(num, 32, &encoded)

	case int64:
		num := reflect.ValueOf(curAppObject).Int()
		encoded = append(encoded, []byte(strconv.FormatInt(0xD3, 16))...)
		numByted(num, 64, &encoded)

	case uint8:
		num := reflect.ValueOf(curAppObject).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCC, 16))...)
		numBytedUnsigned(num, 8, &encoded)

	case uint16:
		num := reflect.ValueOf(curAppObject).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCD, 16))...)
		numBytedUnsigned(num, 16, &encoded)

	case uint32:
		num := reflect.ValueOf(curAppObject).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCE, 16))...)
		numBytedUnsigned(num, 32, &encoded)

	case uint64:
		num := reflect.ValueOf(curAppObject).Uint()
		encoded = append(encoded, []byte(strconv.FormatUint(0xCF, 16))...)
		numBytedUnsigned(num, 64, &encoded)

	case nil:
		encoded = append(encoded, []byte(strconv.FormatUint(0xC0, 16))...)

	case bool:
		if curAppObject {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC3, 16))...)
		} else {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC2, 16))...)
		}

	case float32:
		encoded = append(encoded, []byte(strconv.FormatUint(0xCA, 16))...)
		floatBits := math.Float32bits(curAppObject)
		floatMask := make([]byte, 4)
		binary.BigEndian.PutUint32(floatMask, floatBits)
		hexString := hex.EncodeToString(floatMask)
		encoded = append(encoded, []byte(hexString)...)

	case float64:
		encoded = append(encoded, []byte(strconv.FormatUint(0xCB, 16))...)
		floatBits := math.Float64bits(curAppObject)
		floatMask := make([]byte, 8)
		binary.BigEndian.PutUint64(floatMask, floatBits)
		hexString := hex.EncodeToString(floatMask)
		encoded = append(encoded, []byte(hexString)...)

	case string:
		if len(curAppObject) < 32 {
			if len(curAppObject) < 16 {
				encoded = append(encoded, []byte("a"+strconv.FormatInt(int64(len(curAppObject)), 16))...)
			} else {
				encoded = append(encoded, []byte("b"+strconv.FormatInt(int64(len(curAppObject))-16, 16))...)
			}
		} else if len(curAppObject) < (1 << 8) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xD9, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 8, &encoded)
		} else if len(curAppObject) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDA, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 16, &encoded)
		} else if len(curAppObject) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDB, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		encoded = append(encoded, []byte(curAppObject)...)

	case []byte:
		if len(curAppObject) < (1 << 8) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC4, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 8, &encoded)
		} else if len(curAppObject) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC5, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 16, &encoded)
		} else if len(curAppObject) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xC6, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		encoded = append(encoded, []byte(curAppObject)...)

	case []interface{}:
		if len(curAppObject) < (1 << 4) {
			numBytedUnsigned(uint64(0x90+len(curAppObject)), 8, &encoded)
		} else if len(curAppObject) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDC, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 16, &encoded)
		} else if len(curAppObject) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDD, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		for _, val := range curAppObject {
			elements, err := serializer(val)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, elements...)
		}

	case map[string]interface{}:
		if len(curAppObject) < (1 << 4) {
			numBytedUnsigned(uint64(0x80+len(curAppObject)), 8, &encoded)
		} else if len(curAppObject) < (1 << 16) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDE, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 16, &encoded)
		} else if len(curAppObject) < (1 << 32) {
			encoded = append(encoded, []byte(strconv.FormatUint(0xDF, 16))...)
			numBytedUnsigned(uint64(len(curAppObject)), 32, &encoded)
		} else {
			return nil, fmt.Errorf("unsupported data type")
		}
		for key, val := range curAppObject {
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
	var encodedCopy []byte = *encoded
	return deserializer(&encodedCopy)
}

func deserializer(encoded *[]byte) (interface{}, error) {

	var decoded interface{}

	if len(*encoded) < 2 {
		return nil, fmt.Errorf("unsupported data type")
	}
	var err error
	currObjectType, err := parseUintN((*encoded)[0:2], 8)
	if err != nil {
		return nil, fmt.Errorf("unsupported data type")
	}
	*encoded = (*encoded)[2:]
	switch currObjectType {
	case 0xD0: // int8
		if len(*encoded) < 2 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseIntN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[2:]
	case 0xD1: // int16
		if len(*encoded) < 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseIntN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[4:]
	case 0xD2: // int32
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseIntN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
	case 0xD3: // int64
		if len(*encoded) < 16 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseIntN((*encoded)[:16], 64)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[16:]
	case 0xCC: // uint8
		if len(*encoded) < 2 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseUintN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[2:]
	case 0xCD: // uint16
		if len(*encoded) < 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[4:]
	case 0xCE: // uint32
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
	case 0xCF: // uint64
		if len(*encoded) < 16 {
			return nil, fmt.Errorf("unsupported data type")
		}
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
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseFloatN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[8:]
	case 0xCB: // float64
		if len(*encoded) < 16 {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded, err = parseFloatN((*encoded)[:16], 64)
		if err != nil {
			return nil, err
		}
		*encoded = (*encoded)[16:]
	case 0xD9: // string 8
		if len(*encoded) < 2 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		if len(*encoded) < 2+int(length) {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded = string((*encoded)[2 : 2+length])
		*encoded = (*encoded)[2+length:]
	case 0xDA: // string 16
		if len(*encoded) < 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		if len(*encoded) < 4+int(length) {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded = string((*encoded)[4 : 4+length])
		*encoded = (*encoded)[4+length:]
	case 0xDB: // string 32
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		if len(*encoded) < 8+int(length) {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded = string((*encoded)[8 : 8+length])
		*encoded = (*encoded)[8+length:]
	case 0xC4: // bin 8
		if len(*encoded) < 2 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:2], 8)
		if err != nil {
			return nil, err
		}
		if len(*encoded) < 2+int(length) {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded = (*encoded)[2 : 2+length]
		*encoded = (*encoded)[2+length:]
	case 0xC5: // bin 16
		if len(*encoded) < 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		if len(*encoded) < 4+int(length) {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded = (*encoded)[4 : 4+length]
		*encoded = (*encoded)[4+length:]
	case 0xC6: // bin 32
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		if len(*encoded) < 8+int(length) {
			return nil, fmt.Errorf("unsupported data type")
		}
		decoded = (*encoded)[8 : 8+length]
		*encoded = (*encoded)[8+length:]
	case 0xDC: // array 16
		if len(*encoded) < 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		if len(*encoded) <= 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		*encoded = (*encoded)[4:]
		decoded, err = parseArray(encoded, int(length))
		if err != nil {
			return nil, err
		}
	case 0xDD: // array 32
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		if len(*encoded) <= 8 {
			return nil, fmt.Errorf("array not given elements")
		}
		*encoded = (*encoded)[8:]
		decoded, err = parseArray(encoded, int(length))
		if err != nil {
			return nil, err
		}
	case 0xDE: // map 16
		if len(*encoded) < 4 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:4], 16)
		if err != nil {
			return nil, err
		}
		if len(*encoded) <= 4 {
			return nil, fmt.Errorf("map not given elements")
		}
		*encoded = (*encoded)[4:]
		decoded, err = parseMap(encoded, int(length))
		if err != nil {
			return nil, err
		}
	case 0xDF: // map 32
		if len(*encoded) < 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		length, err := parseUintN((*encoded)[:8], 32)
		if err != nil {
			return nil, err
		}
		if len(*encoded) <= 8 {
			return nil, fmt.Errorf("unsupported data type")
		}
		*encoded = (*encoded)[8:]
		decoded, err = parseMap(encoded, int(length))
		if err != nil {
			return nil, err
		}
	default:
		if currObjectType <= 0x7F {
			decoded = uint8(currObjectType)
		} else if currObjectType <= 0x8f {
			decoded, err = parseMap(encoded, int(currObjectType)-int(0x80))
			if err != nil {
				return nil, err
			}
		} else if currObjectType <= 0x9f {
			decoded, err = parseArray(encoded, int(currObjectType)-int(0x90))
			if err != nil {
				return nil, err
			}
		} else if currObjectType <= 0xbf {
			length := currObjectType - 0xa0
			decoded = string((*encoded)[:length])
			*encoded = (*encoded)[length:]
		} else if currObjectType >= 0xe0 && currObjectType <= 0xff {
			decoded = int8(currObjectType)
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

package core

import (
	"errors"
	"fmt"
)

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	value, _, err := DecodeOne(data)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("empty data")
	}

	switch data[0] {
	case '+':
		return decodeSimpleString(data)
	case '-':
		return decodeError(data)
	case ':':
		return decodeInteger(data)
	case '$':
		return decodeBulkString(data)
	case '*':
		return decodeArray(data)
	default:
		return nil, 0, errors.New("unknown RESP type")
	}
}

func decodeSimpleString(data []byte) (string, int, error) {
	pos := 1

	for ; data[pos] != '\r'; pos++ {

	}
	return string(data[1:pos]), pos + 2, nil
}

func decodeError(data []byte) (string, int, error) {
	return decodeSimpleString(data)
}

func decodeInteger(data []byte) (int64, int, error) {
	pos := 1
	value := int64(0)

	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

func decodeBulkString(data []byte) (string, int, error) {
	pos := 1

	len, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos : pos+len]), pos + len + 2, nil
}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0

	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}

	return 0, 0
}

func decodeArray(data []byte) ([]interface{}, int, error) {
	pos := 1

	count, delta := readLength(data[pos:])
	pos += delta

	values := make([]interface{}, count)

	for i := 0; i < count; i++ {
		value, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		values[i] = value
		pos += delta
	}

	return values, pos, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	ts := value.([]interface{})
	tokens := make([]string, len(ts))
	for i := range tokens {
		tokens[i] = ts[i].(string)
	}

	return tokens, nil
}

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte("+" + v + "\r\n")
		} else {
			return []byte("$" + fmt.Sprintf("%d", len(v)) + "\r\n" + v + "\r\n")
		}
	}
	return []byte{}
}

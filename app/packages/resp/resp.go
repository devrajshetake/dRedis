package resp

import (
	"strconv"
	"strings"
)

func Parse(data []byte) (interface{}, []byte) {
	if len(data) == 0 {
		return nil, data
	}

	switch data[0] {
	case '+':
		return parseSimpleString(data)
	case '-':
		return parseError(data)
	case ':':
		return parseInteger(data)
	case '$':
		return parseBulkString(data)
	case '*':
		return parseArray(data)
	default:
		return nil, data
	}
}

func parseSimpleString(data []byte) (interface{}, []byte) {
	if len(data) < 3 {
		return nil, data
	}

	for i := 1; i < len(data); i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			return string(data[1:i]), data[i+2:]
		}
	}
	return nil, data
}

func parseError(data []byte) (interface{}, []byte) {
	// TODO: Implement this
	return nil, data
}

func parseInteger(data []byte) (interface{}, []byte) {
	// TODO: Implement this
	return nil, data
}

func parseBulkString(data []byte) (interface{}, []byte) {
	if len(data) < 3 || data[0] != '$' {
		return nil, data
	}

	pos := getCrlfIndex(data)
	if pos == -1 {
		return nil, data
	}

	numBytes, err := strconv.Atoi(string(data[1:pos]))
	if err != nil {
		return nil, data
	}

	data = data[pos+2:]

	return string(data[:numBytes]), data[numBytes+2:]
}

func parseArray(data []byte) (interface{}, []byte) {
	// TODO: Implement this
	if len(data) < 3 || data[0] != '*' {
		return nil, data
	}

	pos := getCrlfIndex(data)
	if pos == -1 {
		return nil, data
	}

	numElements, err := strconv.Atoi(string(data[1:pos]))
	if err != nil {
		return nil, data
	}

	data = data[pos+2:]
	elements := make([]interface{}, numElements)

	for i := 0; i < numElements; i++ {
		if len(data) == 0 {
			return nil, data
		}

		switch data[0] {
		case '$':
			element, remaining := parseBulkString(data)
			elements[i] = element
			data = remaining
			// default:
			// 	return nil, data
		}
	}

	return elements, data
}

func getCrlfIndex(data []byte) int {
	return strings.Index(string(data), "\r\n")
}

// Encode
func Encode(data interface{}) []byte {
	switch data := data.(type) {
	case string:
		return encodeSimpleString(data)
	case int:
		return encodeInteger(data)
	case []string:
		return encodeArray(data)
	default:
		return nil
	}
}

func encodeSimpleString(data string) []byte {
	return []byte("+" + data + "\r\n")
}

func encodeInteger(data int) []byte {
	return []byte(":" + strconv.Itoa(data) + "\r\n")
}

func encodeArray(data []string) []byte {
	result := []byte("*" + strconv.Itoa(len(data)) + "\r\n")

	for _, element := range data {
		result = append(result, encodeBulkString(element)...)
	}

	return result
}

func encodeBulkString(data string) []byte {
	return []byte("$" + strconv.Itoa(len(data)) + "\r\n" + data + "\r\n")
}

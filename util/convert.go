package util

import (
	"encoding/json"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func InterfaceToAny[T any](data interface{}) (*T, error) {
	bytes, err := json.Marshal(data)

	var newT = new(T)
	if err != nil {
		return newT, err
	}

	err = json.Unmarshal(bytes, newT)

	return newT, err
}

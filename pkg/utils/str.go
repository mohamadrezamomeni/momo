package utils

import "strconv"

func ConvertToUint16(s string) (uint16, error) {
	num, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(num), nil
}

func ConvertToUint32(s string) (uint32, error) {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

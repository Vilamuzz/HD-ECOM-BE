package helpers

import "strconv"

func ConvertStringToUint64(s string) uint64 {
	var result uint64
	var err error
	if result, err = strconv.ParseUint(s, 10, 64); err != nil {
		return 0
	}
	return result
}

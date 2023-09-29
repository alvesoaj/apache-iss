package utils

import (
	"log"
	"regexp"
	"strconv"
)

var nonNumericRegex = regexp.MustCompile(`[^0-9]+`)

func CastStringToUint8(sVal string, base int) uint8 {
	sVal = RemoveAllNonNumericFromString(sVal)
	nVal, err := strconv.ParseInt(sVal, base, 64)
	if err != nil {
		log.Fatalf("Parsing string to int error: %+v", err)
	}

	return uint8(nVal)
}

func CastInterfaceToUint8(iVal interface{}) uint8 {
	nVal, ok := iVal.(uint8)
	if !ok {
		log.Fatalf("Casting uint8 error, val: %+v", iVal)
	}
	return nVal
}

func RemoveAllNonNumericFromString(sVal string) string {
	return nonNumericRegex.ReplaceAllString(sVal, "")
}

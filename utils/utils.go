package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var nonNumericRegex = regexp.MustCompile(`[^0-9]+`)

func CastStringToUint8(sVal string, base int) uint8 {
	sVal = RemoveAllNonNumericFromString(sVal)
	nVal, err := strconv.ParseInt(sVal, base, 64)
	if err != nil {
		log.Fatalf("Parsing string to uint8 error: %+v", err)
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

func NewTestInput(input string) (*os.File, error) {
	in, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}

	if _, err := io.WriteString(in, input); err != nil {
		in.Close()
		return nil, err
	}

	if _, err := in.Seek(0, io.SeekStart); err != nil {
		in.Close()
		return nil, err
	}

	return in, nil
}

func NewTestOutput() bytes.Buffer {
	var out bytes.Buffer
	return out
}

func ClearOutputForTesting(sVal string) string {
	return strings.Replace(sVal, "> ", "", 3)
}

func CastStringToUint16(sVal string, base int) uint16 {
	sVal = RemoveAllNonNumericFromString(sVal)
	nVal, err := strconv.ParseInt(sVal, base, 64)
	if err != nil {
		log.Fatalf("Parsing string to uint16 error: %+v", err)
	}

	return uint16(nVal)
}

func CastInterfaceToUint16(iVal interface{}) uint16 {
	nVal, ok := iVal.(uint16)
	if !ok {
		log.Fatalf("Casting uint16 error, val: %+v", iVal)
	}
	return nVal
}

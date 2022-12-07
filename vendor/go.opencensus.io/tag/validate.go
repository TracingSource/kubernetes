// Copyright 2017, OpenCensus Authors
//


package tag

import "errors"

const (
	maxKeyLength = 255

	// valid are restricted to US-ASCII subset (range 0x20 (' ') to 0x7e ('~')).
	validKeyValueMin = 32
	validKeyValueMax = 126
)

var (
	errInvalidKeyName = errors.New("invalid key name: only ASCII characters accepted; max length must be 255 characters")
	errInvalidValue   = errors.New("invalid value: only ASCII characters accepted; max length must be 255 characters")
)

func checkKeyName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if len(name) > maxKeyLength {
		return false
	}
	return isASCII(name)
}

func isASCII(s string) bool {
	for _, c := range s {
		if (c < validKeyValueMin) || (c > validKeyValueMax) {
			return false
		}
	}
	return true
}

func checkValue(v string) bool {
	if len(v) > maxKeyLength {
		return false
	}
	return isASCII(v)
}

package flags

import (
	"flag"
	"sort"
	"strings"
)

// UniqueStringsValue wraps a list of unique strings.
// The values are set in order.
type UniqueStringsValue struct {
	Values map[string]struct{}
}

// Set parses a command line set of strings, separated by comma.
// Implements "flag.Value" interface.
// The values are set in order.
func (us *UniqueStringsValue) Set(s string) error {
	us.Values = make(map[string]struct{})
	for _, v := range strings.Split(s, ",") {
		us.Values[v] = struct{}{}
	}
	return nil
}

// String implements "flag.Value" interface.
func (us *UniqueStringsValue) String() string {
	return strings.Join(us.stringSlice(), ",")
}

func (us *UniqueStringsValue) stringSlice() []string {
	ss := make([]string, 0, len(us.Values))
	for v := range us.Values {
		ss = append(ss, v)
	}
	sort.Strings(ss)
	return ss
}

// NewUniqueStringsValue implements string slice as "flag.Value" interface.
// Given value is to be separated by comma.
// The values are set in order.
func NewUniqueStringsValue(s string) (us *UniqueStringsValue) {
	us = &UniqueStringsValue{Values: make(map[string]struct{})}
	if s == "" {
		return us
	}
	if err := us.Set(s); err != nil {
		plog.Panicf("new UniqueStringsValue should never fail: %v", err)
	}
	return us
}

// UniqueStringsFromFlag returns a string slice from the flag.
func UniqueStringsFromFlag(fs *flag.FlagSet, flagName string) []string {
	return (*fs.Lookup(flagName).Value.(*UniqueStringsValue)).stringSlice()
}

// UniqueStringsMapFromFlag returns a map of strings from the flag.
func UniqueStringsMapFromFlag(fs *flag.FlagSet, flagName string) map[string]struct{} {
	return (*fs.Lookup(flagName).Value.(*UniqueStringsValue)).Values
}

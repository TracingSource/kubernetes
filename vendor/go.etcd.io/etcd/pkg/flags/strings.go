package flags

import (
	"flag"
	"sort"
	"strings"
)

// StringsValue wraps "sort.StringSlice".
type StringsValue sort.StringSlice

// Set parses a command line set of strings, separated by comma.
// Implements "flag.Value" interface.
func (ss *StringsValue) Set(s string) error {
	*ss = strings.Split(s, ",")
	return nil
}

// String implements "flag.Value" interface.
func (ss *StringsValue) String() string { return strings.Join(*ss, ",") }

// NewStringsValue implements string slice as "flag.Value" interface.
// Given value is to be separated by comma.
func NewStringsValue(s string) (ss *StringsValue) {
	if s == "" {
		return &StringsValue{}
	}
	ss = new(StringsValue)
	if err := ss.Set(s); err != nil {
		plog.Panicf("new StringsValue should never fail: %v", err)
	}
	return ss
}

// StringsFromFlag returns a string slice from the flag.
func StringsFromFlag(fs *flag.FlagSet, flagName string) []string {
	return []string(*fs.Lookup(flagName).Value.(*StringsValue))
}

package types

import (
	"reflect"
	"strings"
)

var t = map[string]reflect.Type{}

func Add(name string, kind reflect.Type) {
	t[name] = kind
}

type Func func(string) (reflect.Type, bool)

func TypeFunc() Func {
	return func(name string) (reflect.Type, bool) {
		typ, ok := t[name]
		if !ok {
			// The /sdk endpoint does not prefix types with the namespace,
			// but extension endpoints, such as /pbm/sdk do.
			name = strings.TrimPrefix(name, "vim25:")
			typ, ok = t[name]
		}
		return typ, ok
	}
}

package client

import (
	"github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"strconv"
	"unsafe"
)

type customNumberExtension struct {
	jsoniter.DummyExtension
}

func (cne *customNumberExtension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	if typ.String() == "interface {}" {
		return customNumberDecoder{}
	}
	return nil
}

type customNumberDecoder struct {
}

func (customNumberDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		var number jsoniter.Number
		iter.ReadVal(&number)
		i64, err := strconv.ParseInt(string(number), 10, 64)
		if err == nil {
			*(*interface{})(ptr) = i64
			return
		}
		f64, err := strconv.ParseFloat(string(number), 64)
		if err == nil {
			*(*interface{})(ptr) = f64
			return
		}
		iter.ReportError("DecodeNumber", err.Error())
	default:
		*(*interface{})(ptr) = iter.Read()
	}
}

// caseSensitiveJsonIterator returns a jsoniterator API that's configured to be
// case-sensitive when unmarshalling, and otherwise compatible with
// the encoding/json standard library.
func caseSensitiveJsonIterator() jsoniter.API {
	config := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		CaseSensitive:          true,
	}.Froze()
	// Force jsoniter to decode number to interface{} via int64/float64, if possible.
	config.RegisterExtension(&customNumberExtension{})
	return config
}

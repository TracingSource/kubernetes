package runtime

import (
	"encoding/xml"
	"io"
)

// XMLConsumer creates a new XML consumer
func XMLConsumer() Consumer {
	return ConsumerFunc(func(reader io.Reader, data interface{}) error {
		dec := xml.NewDecoder(reader)
		return dec.Decode(data)
	})
}

// XMLProducer creates a new XML producer
func XMLProducer() Producer {
	return ProducerFunc(func(writer io.Writer, data interface{}) error {
		enc := xml.NewEncoder(writer)
		return enc.Encode(data)
	})
}

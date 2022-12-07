package runtime

import (
	"encoding/json"
	"io"
)

// JSONConsumer creates a new JSON consumer
func JSONConsumer() Consumer {
	return ConsumerFunc(func(reader io.Reader, data interface{}) error {
		dec := json.NewDecoder(reader)
		dec.UseNumber() // preserve number formats
		return dec.Decode(data)
	})
}

// JSONProducer creates a new JSON producer
func JSONProducer() Producer {
	return ProducerFunc(func(writer io.Writer, data interface{}) error {
		enc := json.NewEncoder(writer)
		enc.SetEscapeHTML(false)
		return enc.Encode(data)
	})
}

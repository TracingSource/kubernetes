package runtime

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
)

// CSVConsumer creates a new CSV consumer
func CSVConsumer() Consumer {
	return ConsumerFunc(func(reader io.Reader, data interface{}) error {
		if reader == nil {
			return errors.New("CSVConsumer requires a reader")
		}

		csvReader := csv.NewReader(reader)
		writer, ok := data.(io.Writer)
		if !ok {
			return errors.New("data type must be io.Writer")
		}
		csvWriter := csv.NewWriter(writer)
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}
		for _, r := range records {
			if err := csvWriter.Write(r); err != nil {
				return err
			}
		}
		csvWriter.Flush()
		return nil
	})
}

// CSVProducer creates a new CSV producer
func CSVProducer() Producer {
	return ProducerFunc(func(writer io.Writer, data interface{}) error {
		if writer == nil {
			return errors.New("CSVProducer requires a writer")
		}

		dataBytes, ok := data.([]byte)
		if !ok {
			return errors.New("data type must be byte array")
		}

		csvReader := csv.NewReader(bytes.NewBuffer(dataBytes))
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}
		csvWriter := csv.NewWriter(writer)
		for _, r := range records {
			if err := csvWriter.Write(r); err != nil {
				return err
			}
		}
		csvWriter.Flush()
		return nil
	})
}

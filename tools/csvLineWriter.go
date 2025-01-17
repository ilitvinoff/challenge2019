package tools

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CsvWriter interface {
	WriteLineToCsv(rowChan <-chan []string, errChan chan<- error)
}

type CsvWriteConfig struct {
	DestinationFilepath string
	createDestination   func(string) (io.WriteCloser, error)
}

func NewCsvWriterConfig(destinationFilepath string) *CsvWriteConfig {
	return &CsvWriteConfig{destinationFilepath, createDestination}
}

// WriteLineToCsv writes row to destination csv file. If dest file exists, the new one will be created to replace existent.
// Row is received from chan []string, if any error acquired, send it to errChan and stops writing.
func (cfg *CsvWriteConfig) WriteLineToCsv(rowChan <-chan []string, errChan chan<- error) {
	defer close(errChan)
	f, err := cfg.createDestination(cfg.DestinationFilepath)
	if err != nil {
		errChan <- fmt.Errorf("can't create destination file: {%s}; err: %s", cfg.DestinationFilepath, err)
		return
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)

	for row := range rowChan {
		err := csvWriter.Write(row)
		if err != nil {
			errChan <- fmt.Errorf("can't write row {%s} to csv: %s", row, err)
			return
		}

	}

	csvWriter.Flush()
}

func createDestination(dest string) (io.WriteCloser, error) {
	f, err := os.Create(dest)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Package csv provides functionality for parsing and processing CSV files.
// It includes functions to read CSV files and convert the records to maps.
package csv

import (
	"encoding/csv"
	"io"
	"os"
)

// Parser provides CSV parsing functionality
type Parser struct{}

// NewParser creates a new CSV parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads a CSV file and returns records and headers
func (p *Parser) Parse(filePath string) ([][]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	// Read remaining records
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		records = append(records, record)
	}

	return records, headers, nil
}

// MapRecords converts CSV records to maps using headers as keys
func (p *Parser) MapRecords(records [][]string, headers []string) []map[string]string {
	result := make([]map[string]string, 0, len(records))

	for _, record := range records {
		data := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				data[header] = record[i]
			}
		}
		result = append(result, data)
	}

	return result
}

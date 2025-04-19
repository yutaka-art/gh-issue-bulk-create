// Package csv provides functionality for parsing and processing CSV files.
// It includes functions to read CSV files and convert the records to maps.
package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
		if err == io.EOF {
			return nil, nil, errors.New("CSV file is empty")
		}
		return nil, nil, err
	}

	// Validate that headers are not empty
	if len(headers) == 0 {
		return nil, nil, errors.New("CSV file has no headers")
	}

	// Validate that headers do not contain empty strings
	for i, header := range headers {
		if header == "" {
			return nil, nil, errors.New("empty header found at column " + string(rune('A'+i)))
		}
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

// ValidateHeadersAgainstTemplate validates that CSV headers match the variables in the template
func (p *Parser) ValidateHeadersAgainstTemplate(headers []string, templateVars []string) ([]string, error) {
	if len(headers) == 0 {
		return nil, errors.New("CSV has no headers")
	}

	// Create a map of template variables for quick lookup
	templateVarMap := make(map[string]bool)
	for _, v := range templateVars {
		templateVarMap[v] = true
	}

	// Check if each header exists in the template variables
	var missingVars []string
	var missingHeaders []string

	// Check for CSV headers that don't exist in template
	for _, header := range headers {
		if !templateVarMap[header] {
			missingVars = append(missingVars, header)
		}
	}

	// Check for template variables that don't exist in CSV headers
	headerMap := make(map[string]bool)
	for _, h := range headers {
		headerMap[h] = true
	}
	for _, v := range templateVars {
		if !headerMap[v] {
			missingHeaders = append(missingHeaders, v)
		}
	}

	var warnings []string

	if len(missingVars) > 0 {
		warnings = append(warnings, fmt.Sprintf("Warning: The following CSV headers are not used in the template: %s", strings.Join(missingVars, ", ")))
	}

	if len(missingHeaders) > 0 {
		warnings = append(warnings, fmt.Sprintf("Warning: The following template variables are missing from CSV headers: %s", strings.Join(missingHeaders, ", ")))
	}

	return warnings, nil
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

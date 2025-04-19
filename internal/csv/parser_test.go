package csv

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	// Create a temporary CSV file for testing
	tmpFile, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data to the CSV file
	writer := csv.NewWriter(tmpFile)
	testData := [][]string{
		{"header1", "header2", "header3"},
		{"value1", "value2", "value3"},
		{"value4", "value5", "value6"},
	}
	if err := writer.WriteAll(testData); err != nil {
		t.Fatalf("Failed to write to CSV file: %v", err)
	}
	writer.Flush()
	tmpFile.Close()

	// Test the Parse function
	parser := NewParser()
	records, headers, err := parser.Parse(tmpFile.Name())
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check headers
	expectedHeaders := testData[0]
	if !reflect.DeepEqual(headers, expectedHeaders) {
		t.Errorf("Expected headers %v, got %v", expectedHeaders, headers)
	}

	// Check records
	expectedRecords := testData[1:]
	if !reflect.DeepEqual(records, expectedRecords) {
		t.Errorf("Expected records %v, got %v", expectedRecords, records)
	}
}

func TestMapRecords(t *testing.T) {
	// Test data
	headers := []string{"header1", "header2", "header3"}
	records := [][]string{
		{"value1", "value2", "value3"},
		{"value4", "value5", "value6"},
	}

	// Expected mapped records
	expected := []map[string]string{
		{
			"header1": "value1",
			"header2": "value2",
			"header3": "value3",
		},
		{
			"header1": "value4",
			"header2": "value5",
			"header3": "value6",
		},
	}

	// Test the MapRecords function
	parser := NewParser()
	result := parser.MapRecords(records, headers)

	// Check result
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected mapped records %v, got %v", expected, result)
	}
}

func TestMapRecordsWithMissingValues(t *testing.T) {
	// Test data with missing values
	headers := []string{"header1", "header2", "header3"}
	records := [][]string{
		{"value1", "value2"}, // Missing value for header3
		{"value4"},           // Missing values for header2 and header3
	}

	// Expected mapped records
	expected := []map[string]string{
		{
			"header1": "value1",
			"header2": "value2",
			// header3 is missing
		},
		{
			"header1": "value4",
			// header2 and header3 are missing
		},
	}

	// Test the MapRecords function
	parser := NewParser()
	result := parser.MapRecords(records, headers)

	// Check result
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected mapped records %v, got %v", expected, result)
	}
}

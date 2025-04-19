package csv

import (
	"encoding/csv"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestValidateHeadersAgainstTemplate(t *testing.T) {
	testCases := []struct {
		name          string
		headers       []string
		templateVars  []string
		expectWarning bool
		warningSubstr string
	}{
		{
			name:          "All headers match template variables",
			headers:       []string{"title", "description", "steps"},
			templateVars:  []string{"title", "description", "steps"},
			expectWarning: false,
		},
		{
			name:          "Missing headers",
			headers:       []string{"title", "description"},
			templateVars:  []string{"title", "description", "steps"},
			expectWarning: true,
			warningSubstr: "missing from CSV headers",
		},
		{
			name:          "Extra headers",
			headers:       []string{"title", "description", "steps", "extra"},
			templateVars:  []string{"title", "description", "steps"},
			expectWarning: true,
			warningSubstr: "not used in the template",
		},
		{
			name:          "Both missing and extra headers",
			headers:       []string{"title", "extra1", "extra2"},
			templateVars:  []string{"title", "description", "steps"},
			expectWarning: true,
			warningSubstr: "missing from CSV headers",
		},
		{
			name:          "Empty headers list",
			headers:       []string{},
			templateVars:  []string{"title", "description"},
			expectWarning: false, // This will return an error, not a warning
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser()
			warnings, err := parser.ValidateHeadersAgainstTemplate(tc.headers, tc.templateVars)

			if len(tc.headers) == 0 {
				if err == nil {
					t.Error("Expected error for empty headers, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tc.expectWarning {
				if len(warnings) == 0 {
					t.Error("Expected warnings, but got none")
				} else {
					foundExpectedWarning := false
					for _, warning := range warnings {
						if strings.Contains(warning, tc.warningSubstr) {
							foundExpectedWarning = true
							break
						}
					}
					if !foundExpectedWarning {
						t.Errorf("Expected warning containing '%s', but not found in %v", tc.warningSubstr, warnings)
					}
				}
			} else {
				if len(warnings) > 0 {
					t.Errorf("Did not expect warnings, but got: %v", warnings)
				}
			}
		})
	}
}

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

func TestParseEmptyFile(t *testing.T) {
	// Create an empty temporary CSV file for testing
	tmpFile, err := os.CreateTemp("", "empty-test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test the Parse function with empty file
	parser := NewParser()
	_, _, err = parser.Parse(tmpFile.Name())

	// Should return error for empty file
	if err == nil {
		t.Error("Expected error for empty file, but got nil")
	}

	if !strings.Contains(err.Error(), "CSV file is empty") {
		t.Errorf("Expected error message to contain 'CSV file is empty', got: %v", err)
	}
}

func TestParseEmptyHeaders(t *testing.T) {
	// Create a temporary CSV file with empty headers
	tmpFile, err := os.CreateTemp("", "empty-headers-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write a file with an empty header column
	content := "header1,,header3\nvalue1,value2,value3"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to CSV file: %v", err)
	}
	tmpFile.Close()

	// Test the Parse function with empty header
	parser := NewParser()
	_, _, err = parser.Parse(tmpFile.Name())

	// Should return error for empty header
	if err == nil {
		t.Error("Expected error for empty header, but got nil")
	}

	if !strings.Contains(err.Error(), "empty header found at column") {
		t.Errorf("Expected error message to contain 'empty header found at column', got: %v", err)
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

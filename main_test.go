package main

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
)

func TestReadCSV(t *testing.T) {
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

	// Test the readCSV function
	records, headers, err := readCSV(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSV failed: %v", err)
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

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]string
		expected string
		wantErr  bool
	}{
		{
			name:     "Basic template",
			template: "Hello {{name}}!",
			data:     map[string]string{"name": "World"},
			expected: "Hello World!",
			wantErr:  false,
		},
		{
			name:     "Multiple variables",
			template: "{{greeting}} {{name}}!",
			data:     map[string]string{"greeting": "Hello", "name": "World"},
			expected: "Hello World!",
			wantErr:  false,
		},
		{
			name:     "Missing variable",
			template: "Hello {{name}}!",
			data:     map[string]string{},
			expected: "Hello !",
			wantErr:  false,
		},
		{
			name:     "Invalid template",
			template: "Hello {{name!",
			data:     map[string]string{"name": "World"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderTemplate(tt.template, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("renderTemplate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseIssueTemplate(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected IssueData
		wantErr  bool
	}{
		{
			name: "Valid issue template",
			content: `---
title: Test Issue
labels: bug, enhancement
assignees: user1, user2
---
This is the body of the issue.`,
			expected: IssueData{
				Title:     "Test Issue",
				Labels:    []string{"bug", "enhancement"},
				Assignees: []string{"user1", "user2"},
				Body:      "This is the body of the issue.",
			},
			wantErr: false,
		},
		{
			name: "Array labels and assignees",
			content: `---
title: Test Issue
labels: 
  - bug
  - enhancement
assignees: 
  - user1
  - user2
---
This is the body of the issue.`,
			expected: IssueData{
				Title:     "Test Issue",
				Labels:    []string{"bug", "enhancement"},
				Assignees: []string{"user1", "user2"},
				Body:      "This is the body of the issue.",
			},
			wantErr: false,
		},
		{
			name:     "Missing front matter delimiter",
			content:  "title: Test Issue\nThis is the body of the issue.",
			expected: IssueData{},
			wantErr:  true,
		},
		{
			name: "Incomplete front matter",
			content: `---
title: Test Issue
---
This is the body of the issue.`,
			expected: IssueData{
				Title: "Test Issue",
				Body:  "This is the body of the issue.",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseIssueTemplate(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIssueTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Title != tt.expected.Title {
					t.Errorf("Title = %v, want %v", result.Title, tt.expected.Title)
				}
				if !reflect.DeepEqual(result.Labels, tt.expected.Labels) {
					t.Errorf("Labels = %v, want %v", result.Labels, tt.expected.Labels)
				}
				if !reflect.DeepEqual(result.Assignees, tt.expected.Assignees) {
					t.Errorf("Assignees = %v, want %v", result.Assignees, tt.expected.Assignees)
				}
				if result.Body != tt.expected.Body {
					t.Errorf("Body = %v, want %v", result.Body, tt.expected.Body)
				}
			}
		})
	}
}

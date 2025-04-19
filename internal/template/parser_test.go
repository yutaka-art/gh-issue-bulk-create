package template

import (
	"reflect"
	"testing"
)

func TestParseIssueTemplate(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		content        string
		expectedTitle  string
		expectedBody   string
		expectedLabels []string
		expectedError  bool
	}{
		{
			name: "Valid issue template",
			content: `---
title: "Test Issue"
labels: bug, enhancement
assignees: user1, user2
---
This is the body of the issue.
With multiple lines.`,
			expectedTitle:  "Test Issue",
			expectedBody:   "This is the body of the issue.\nWith multiple lines.",
			expectedLabels: []string{"bug", "enhancement"},
			expectedError:  false,
		},
		{
			name:          "No front matter",
			content:       "This is just content without front matter.",
			expectedError: true,
		},
		{
			name: "Invalid front matter format",
			content: `---
This is not valid YAML
---
Content`,
			expectedError: true,
		},
		{
			name: "Missing closing delimiter",
			content: `---
title: "Test"
Content without closing front matter delimiter`,
			expectedError: true,
		},
		{
			name: "Array labels",
			content: `---
title: "Test Issue"
labels:
  - bug
  - enhancement
---
Body content`,
			expectedTitle:  "Test Issue",
			expectedBody:   "Body content",
			expectedLabels: []string{"bug", "enhancement"},
			expectedError:  false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser()
			issue, err := parser.ParseIssueTemplate(tc.content)

			// Check error expectation
			if tc.expectedError && err == nil {
				t.Error("Expected error, got nil")
				return
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			// Skip further checks if we expected an error
			if tc.expectedError {
				return
			}

			// Check issue properties
			if issue.Title != tc.expectedTitle {
				t.Errorf("Expected title '%s', got '%s'", tc.expectedTitle, issue.Title)
			}

			if issue.Body != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, issue.Body)
			}

			if !reflect.DeepEqual(issue.Labels, tc.expectedLabels) {
				t.Errorf("Expected labels %v, got %v", tc.expectedLabels, issue.Labels)
			}
		})
	}
}

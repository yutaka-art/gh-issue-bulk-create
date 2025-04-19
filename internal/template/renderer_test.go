package template

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestExtractVariables(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		expected []string
	}{
		{
			name:     "Simple template",
			template: "Hello, {{name}}!",
			expected: []string{"name"},
		},
		{
			name:     "Multiple variables",
			template: "{{greeting}}, {{name}}!",
			expected: []string{"greeting", "name"},
		},
		{
			name:     "Duplicate variables",
			template: "{{name}}, {{name}}!",
			expected: []string{"name"},
		},
		{
			name:     "Variables with whitespace",
			template: "{{ name }}, {{ greeting }}!",
			expected: []string{"name", "greeting"},
		},
		{
			name: "Markdown template with frontmatter",
			template: `---
title: "{{title}}"
labels: {{label1}}, {{label2}}
---
## Description
{{description}}
## Steps
{{steps}}`,
			expected: []string{"title", "label1", "label2", "description", "steps"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			renderer := NewRenderer()
			result := renderer.ExtractVariables(tc.template)

			// Sort both slices for comparison
			sort.Strings(result)
			sort.Strings(tc.expected)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected variables %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestRender(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		template string
		data     map[string]string
		expected string
	}{
		{
			name:     "Simple template",
			template: "Hello, {{name}}!",
			data:     map[string]string{"name": "World"},
			expected: "Hello, World!",
		},
		{
			name:     "Multiple variables",
			template: "{{greeting}}, {{name}}!",
			data:     map[string]string{"greeting": "Hello", "name": "World"},
			expected: "Hello, World!",
		},
		{
			name:     "Missing variable",
			template: "Hello, {{name}}! Today is {{day}}.",
			data:     map[string]string{"name": "World"},
			expected: "Hello, World! Today is .",
		},
		{
			name:     "Multiline template",
			template: "Title: {{title}}\nDescription: {{description}}",
			data:     map[string]string{"title": "Test", "description": "This is a test"},
			expected: "Title: Test\nDescription: This is a test",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			renderer := NewRenderer()
			result, err := renderer.Render(tc.template, tc.data)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestRenderWithInvalidTemplate(t *testing.T) {
	// Test invalid template
	invalidTemplate := "Hello, {{name!" // Missing closing brace
	data := map[string]string{"name": "World"}

	renderer := NewRenderer()
	_, err := renderer.Render(invalidTemplate, data)
	if err == nil {
		t.Error("Expected error with invalid template, got nil")
	}

	// Check if the error message contains an error indication, without being specific about the exact message
	if err != nil && !strings.Contains(err.Error(), "template") && !strings.Contains(err.Error(), "error") {
		t.Errorf("Expected error to contain template error information, got '%s'", err.Error())
	}
}

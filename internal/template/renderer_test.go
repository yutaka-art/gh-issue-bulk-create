package template

import (
	"strings"
	"testing"
)

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

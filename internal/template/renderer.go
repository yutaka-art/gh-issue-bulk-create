// Package template provides functionality for parsing and rendering markdown templates.
// It includes both parsing of front matter in markdown files and rendering templates with data.
package template

import (
	"regexp"
	"strings"
	"text/template"
)

// Renderer provides template rendering functionality
type Renderer struct{}

// NewRenderer creates a new template renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Render processes a template with the provided data
func (r *Renderer) Render(tmplContent string, data map[string]string) (string, error) {
	// Extract all variable names from the template using a regexp
	re := regexp.MustCompile(`{{([^}]+)}}`)
	matches := re.FindAllStringSubmatch(tmplContent, -1)

	// Replace variables with Go template syntax
	for _, match := range matches {
		if len(match) > 1 {
			varName := strings.TrimSpace(match[1])
			// Replace {{varName}} with {{$.varName}} but only if the variable is in data
			// This prevents errors when the template contains variables not in the data map
			if _, exists := data[varName]; exists {
				tmplContent = strings.ReplaceAll(tmplContent, "{{"+varName+"}}", "{{$."+varName+"}}")
			} else {
				// For non-existent variables, replace with empty string
				tmplContent = strings.ReplaceAll(tmplContent, "{{"+varName+"}}", "")
			}
		}
	}

	tmpl, err := template.New("issue").Delims("{{", "}}").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

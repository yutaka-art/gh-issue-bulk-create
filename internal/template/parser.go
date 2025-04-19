// Package template provides functionality for parsing and rendering markdown templates.
// It includes both parsing of front matter in markdown files and rendering templates with data.
package template

import (
	"fmt"
	"strings"

	"github.com/ntsk/gh-issue-bulk-create/pkg/models"
	"gopkg.in/yaml.v3"
)

// Parser provides markdown parsing functionality
type Parser struct{}

// NewParser creates a new markdown parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseIssueTemplate parses a markdown template with front matter
// and returns an Issue model
func (p *Parser) ParseIssueTemplate(content string) (*models.Issue, error) {
	// Check if the content contains front matter
	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("content does not start with front matter delimiter '---'")
	}

	// Split content into front matter and body
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid front matter format")
	}

	frontMatter := parts[1]
	body := strings.TrimSpace(parts[2])

	// Parse front matter as YAML
	metadata := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(frontMatter), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %v", err)
	}

	// Extract metadata and create the Issue
	var issue models.Issue
	issue.Body = body

	// Extract title
	if title, ok := metadata["title"].(string); ok {
		issue.Title = title
	}

	// Extract labels
	if labels, ok := metadata["labels"].(string); ok {
		labelList := strings.Split(labels, ",")
		for i, label := range labelList {
			labelList[i] = strings.TrimSpace(label)
		}
		issue.Labels = labelList
	} else if labelsArray, ok := metadata["labels"].([]interface{}); ok {
		for _, label := range labelsArray {
			if labelStr, ok := label.(string); ok {
				issue.Labels = append(issue.Labels, labelStr)
			}
		}
	}

	// Extract assignees
	if assignees, ok := metadata["assignees"].(string); ok {
		assigneeList := strings.Split(assignees, ",")
		for i, assignee := range assigneeList {
			assigneeList[i] = strings.TrimSpace(assignee)
		}
		issue.Assignees = assigneeList
	} else if assigneesArray, ok := metadata["assignees"].([]interface{}); ok {
		for _, assignee := range assigneesArray {
			if assigneeStr, ok := assignee.(string); ok {
				issue.Assignees = append(issue.Assignees, assigneeStr)
			}
		}
	}

	// Extract milestone
	if milestone, ok := metadata["milestone"].(string); ok {
		issue.Milestone = milestone
	}

	return &issue, nil
}

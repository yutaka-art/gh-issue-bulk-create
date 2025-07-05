// Package models provides data structures used throughout the application.
// It includes models for GitHub issues and related entities.
package models

// Issue represents a GitHub issue with its metadata
type Issue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Milestone string   `json:"milestone,omitempty"`
}

// NewIssue creates a new Issue with the given title and body
func NewIssue(title, body string) *Issue {
	return &Issue{
		Title: title,
		Body:  body,
	}
}

// WithLabels adds labels to the issue
func (i *Issue) WithLabels(labels []string) *Issue {
	i.Labels = labels
	return i
}

// WithAssignees adds assignees to the issue
func (i *Issue) WithAssignees(assignees []string) *Issue {
	i.Assignees = assignees
	return i
}

// WithMilestone adds a milestone to the issue
func (i *Issue) WithMilestone(milestone string) *Issue {
	i.Milestone = milestone
	return i
}

// IssueResponse represents a GitHub API response when creating an issue
type IssueResponse struct {
	Number int    `json:"number"`
	URL    string `json:"html_url"`
}

// RateLimit represents GitHub API rate limit information
type RateLimit struct {
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"`
}

// RateLimitResponse represents GitHub API rate limit response
type RateLimitResponse struct {
	Rate RateLimit `json:"rate"`
}

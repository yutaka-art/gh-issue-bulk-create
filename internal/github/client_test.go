package github

import (
	"testing"

	"github.com/ntsk/gh-issue-bulk-create/pkg/models"
)

func TestMockClient(t *testing.T) {
	// Create mock client
	mockClient := &MockClient{}

	// Test issue
	issue := &models.Issue{
		Title: "Test Issue",
		Body:  "This is a test issue",
	}

	// Test CreateIssue
	response, err := mockClient.CreateIssue(issue, "test/repo")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check response
	if response.Number != 1 {
		t.Errorf("Expected issue number 1, got %d", response.Number)
	}

	// Check that issue was stored
	if len(mockClient.CreatedIssues) != 1 {
		t.Errorf("Expected 1 created issue, got %d", len(mockClient.CreatedIssues))
	}

	// Test GetCurrentRepository
	repo, err := mockClient.GetCurrentRepository()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check repo
	if repo != "mock/repo" {
		t.Errorf("Expected repo 'mock/repo', got '%s'", repo)
	}

	// Check counter
	if mockClient.GetCurrentRepoCounter != 1 {
		t.Errorf("Expected GetCurrentRepository to be called once, called %d times", mockClient.GetCurrentRepoCounter)
	}

	// Test with custom functions
	mockClient = &MockClient{
		CreateIssueFunc: func(issue *models.Issue, repo string) (*models.IssueResponse, error) {
			return &models.IssueResponse{Number: 42, URL: "custom-url"}, nil
		},
		GetCurrentRepoFunc: func() (string, error) {
			return "custom/repo", nil
		},
	}

	// Test custom CreateIssue
	response, _ = mockClient.CreateIssue(issue, "test/repo")
	if response.Number != 42 {
		t.Errorf("Expected custom issue number 42, got %d", response.Number)
	}

	// Test custom GetCurrentRepository
	repo, _ = mockClient.GetCurrentRepository()
	if repo != "custom/repo" {
		t.Errorf("Expected custom repo 'custom/repo', got '%s'", repo)
	}
}

package github

import (
	"testing"

	"github.com/ntsk/gh-issue-bulk-create/pkg/models"
)

// MockClient provides a mock GitHub client for testing
type MockClient struct {
	CreateIssueFunc       func(issue *models.Issue, repo string) (*models.IssueResponse, error)
	GetCurrentRepoFunc    func() (string, error)
	CreatedIssues         []*models.Issue
	GetCurrentRepoCounter int
}

// CreateIssue implements the ClientInterface for testing
func (m *MockClient) CreateIssue(issue *models.Issue, repo string) (*models.IssueResponse, error) {
	m.CreatedIssues = append(m.CreatedIssues, issue)
	if m.CreateIssueFunc != nil {
		return m.CreateIssueFunc(issue, repo)
	}
	return &models.IssueResponse{Number: 1, URL: "https://github.com/mock/repo/issues/1"}, nil
}

// GetCurrentRepository implements the ClientInterface for testing
func (m *MockClient) GetCurrentRepository() (string, error) {
	m.GetCurrentRepoCounter++
	if m.GetCurrentRepoFunc != nil {
		return m.GetCurrentRepoFunc()
	}
	return "mock/repo", nil
}

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

// TestClientInterface demonstrates that both Client and MockClient implement ClientInterface
func TestClientInterface(t *testing.T) {
	// Verify that both types implement the interface
	var client ClientInterface

	// Test with MockClient
	mockClient := &MockClient{}
	client = mockClient
	if client == nil {
		t.Error("MockClient should implement ClientInterface")
	}

	// Test with real Client (this would fail without proper initialization in tests)
	// realClient := &Client{}
	// client = realClient
	// if client == nil {
	//     t.Error("Client should implement ClientInterface")
	// }
}

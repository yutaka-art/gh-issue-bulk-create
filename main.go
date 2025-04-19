package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/api"
	"gopkg.in/yaml.v3"
)

func main() {
	// Define command line arguments
	templateFile := flag.String("template", "", "Path to the template markdown file")
	csvFile := flag.String("csv", "", "Path to the CSV file containing data")
	dryRun := flag.Bool("dry-run", false, "Only show the content of issues without creating them")
	repo := flag.String("repo", "", "Target repository in the format of owner/repo (default: current repository)")

	flag.Parse()

	// Check arguments
	if *templateFile == "" || *csvFile == "" {
		fmt.Println("Error: Both template file and CSV file must be specified")
		flag.Usage()
		os.Exit(1)
	}

	// Read template file
	tmplContent, err := os.ReadFile(*templateFile)
	if err != nil {
		fmt.Printf("Failed to read template file: %v\n", err)
		os.Exit(1)
	}

	// Read CSV file
	records, headers, err := readCSV(*csvFile)
	if err != nil {
		fmt.Printf("Failed to read CSV file: %v\n", err)
		os.Exit(1)
	}

	// Initialize GitHub API client
	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Printf("Failed to initialize GitHub API client: %v\n", err)
		os.Exit(1)
	}

	// Determine repository
	targetRepo := *repo
	if targetRepo == "" {
		// If not specified as a flag, try to get from current directory
		targetRepo, err = getCurrentRepository()
		if err != nil {
			fmt.Printf("Failed to determine repository: %v\n", err)
			fmt.Println("Please specify the repository using --repo option or run in a git repository")
			os.Exit(1)
		}
	}

	fmt.Printf("Target repository: %s\n", targetRepo)

	// Process template and create issues
	for _, record := range records {
		// Convert record data to map
		data := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				data[header] = record[i]
			}
		}

		// Process template
		processedContent, err := renderTemplate(string(tmplContent), data)
		if err != nil {
			fmt.Printf("Failed to process template: %v\n", err)
			continue
		}

		// Parse issue template to separate metadata and body
		issueData, err := parseIssueTemplate(processedContent)
		if err != nil {
			fmt.Printf("Failed to parse issue template: %v\n", err)
			continue
		}

		if *dryRun {
			// Dry run: Show issue content
			fmt.Println("==== Issue Content ====")
			fmt.Printf("Title: %s\n", issueData.Title)
			fmt.Printf("Labels: %v\n", issueData.Labels)
			fmt.Printf("Assignees: %v\n", issueData.Assignees)
			fmt.Printf("Body:\n%s\n", issueData.Body)
			fmt.Println("=====================")
		} else {
			// Create issue
			err = createIssue(client, issueData, targetRepo)
			if err != nil {
				fmt.Printf("Failed to create issue: %v\n", err)
			} else {
				fmt.Printf("Successfully created issue: %s\n", issueData.Title)
			}
		}
	}
}

// IssueData stores issue metadata
type IssueData struct {
	Title     string
	Body      string
	Labels    []string
	Assignees []string
	Milestone string
}

// Read CSV file and return records and headers
func readCSV(filePath string) ([][]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	// Read remaining records
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		records = append(records, record)
	}

	return records, headers, nil
}

// Render template string using data
func renderTemplate(tmplContent string, data map[string]string) (string, error) {
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

// Parse issue template and return IssueData
func parseIssueTemplate(content string) (IssueData, error) {
	var issueData IssueData

	// Check if the content contains front matter
	if !strings.HasPrefix(content, "---") {
		return issueData, fmt.Errorf("content does not start with front matter delimiter '---'")
	}

	// Split content into front matter and body
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return issueData, fmt.Errorf("invalid front matter format")
	}

	frontMatter := parts[1]
	body := strings.TrimSpace(parts[2])

	// Parse front matter as YAML
	metadata := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(frontMatter), &metadata)
	if err != nil {
		return issueData, fmt.Errorf("failed to parse front matter: %v", err)
	}

	// Extract metadata
	if title, ok := metadata["title"].(string); ok {
		issueData.Title = title
	}

	if labels, ok := metadata["labels"].(string); ok {
		labelList := strings.Split(labels, ",")
		for i, label := range labelList {
			labelList[i] = strings.TrimSpace(label)
		}
		issueData.Labels = labelList
	} else if labelsArray, ok := metadata["labels"].([]interface{}); ok {
		for _, label := range labelsArray {
			if labelStr, ok := label.(string); ok {
				issueData.Labels = append(issueData.Labels, labelStr)
			}
		}
	}

	if assignees, ok := metadata["assignees"].(string); ok {
		assigneeList := strings.Split(assignees, ",")
		for i, assignee := range assigneeList {
			assigneeList[i] = strings.TrimSpace(assignee)
		}
		issueData.Assignees = assigneeList
	} else if assigneesArray, ok := metadata["assignees"].([]interface{}); ok {
		for _, assignee := range assigneesArray {
			if assigneeStr, ok := assignee.(string); ok {
				issueData.Assignees = append(issueData.Assignees, assigneeStr)
			}
		}
	}

	if milestone, ok := metadata["milestone"].(string); ok {
		issueData.Milestone = milestone
	}

	// Set body
	issueData.Body = body

	return issueData, nil
}

// Create issue using GitHub API
func createIssue(client *api.RESTClient, issueData IssueData, repo string) error {
	// GitHub API response structure
	response := struct {
		Number int    `json:"number"`
		URL    string `json:"html_url"`
	}{}

	// Build request body
	requestBody := map[string]interface{}{
		"title":     issueData.Title,
		"body":      issueData.Body,
		"labels":    issueData.Labels,
		"assignees": issueData.Assignees,
	}

	if issueData.Milestone != "" {
		requestBody["milestone"] = issueData.Milestone
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send POST request
	path := fmt.Sprintf("repos/%s/issues", repo)
	err = client.Post(path, bytes.NewReader(jsonData), &response)
	if err != nil {
		return err
	}

	fmt.Printf("Issue #%d created: %s\n", response.Number, response.URL)
	return nil
}

// Get repository information from current directory
func getCurrentRepository() (string, error) {
	// Get repository information using gh command
	type RepoInfo struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	}

	// Use gh command to get repository information
	output, stderr, err := gh.Exec("repo", "view", "--json", "owner,name")
	if err != nil {
		return "", fmt.Errorf("failed to get repository information: %s - %s", err, stderr.String())
	}

	// Parse output as JSON
	var info RepoInfo
	if err := json.Unmarshal(output.Bytes(), &info); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", info.Owner.Login, info.Name), nil
}

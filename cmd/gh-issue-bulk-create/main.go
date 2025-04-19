package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ntsk/gh-issue-bulk-create/internal/csv"
	"github.com/ntsk/gh-issue-bulk-create/internal/github"
	"github.com/ntsk/gh-issue-bulk-create/internal/template"
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

	// Initialize components
	csvParser := csv.NewParser()
	templateRenderer := template.NewRenderer()
	templateParser := template.NewParser()

	// Initialize GitHub client
	githubClient, err := github.NewClient()
	if err != nil {
		fmt.Printf("Failed to initialize GitHub API client: %v\n", err)
		os.Exit(1)
	}

	// Read template file
	tmplContent, err := os.ReadFile(*templateFile)
	if err != nil {
		fmt.Printf("Failed to read template file: %v\n", err)
		os.Exit(1)
	}

	// Read CSV file
	records, headers, err := csvParser.Parse(*csvFile)
	if err != nil {
		fmt.Printf("Failed to read CSV file: %v\n", err)
		os.Exit(1)
	}

	// Map records to data maps
	dataMaps := csvParser.MapRecords(records, headers)

	// Determine repository
	targetRepo := *repo
	if targetRepo == "" {
		// If not specified as a flag, try to get from current directory
		targetRepo, err = githubClient.GetCurrentRepository()
		if err != nil {
			fmt.Printf("Failed to determine repository: %v\n", err)
			fmt.Println("Please specify the repository using --repo option or run in a git repository")
			os.Exit(1)
		}
	}

	fmt.Printf("Target repository: %s\n", targetRepo)

	// Process template and create issues
	for _, data := range dataMaps {
		// Render template with data
		processedContent, err := templateRenderer.Render(string(tmplContent), data)
		if err != nil {
			fmt.Printf("Failed to process template: %v\n", err)
			continue
		}

		// Parse issue template to get issue data
		issue, err := templateParser.ParseIssueTemplate(processedContent)
		if err != nil {
			fmt.Printf("Failed to parse issue template: %v\n", err)
			continue
		}

		if *dryRun {
			// Dry run: Show issue content
			fmt.Println("==== Issue Content ====")
			fmt.Printf("Title: %s\n", issue.Title)
			fmt.Printf("Labels: %v\n", issue.Labels)
			fmt.Printf("Assignees: %v\n", issue.Assignees)
			fmt.Printf("Body:\n%s\n", issue.Body)
			fmt.Println("=====================")
		} else {
			// Create issue
			response, err := githubClient.CreateIssue(issue, targetRepo)
			if err != nil {
				fmt.Printf("Failed to create issue: %v\n", err)
			} else {
				fmt.Printf("Issue #%d created: %s\n", response.Number, response.URL)
			}
		}
	}
}

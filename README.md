# gh-issue-bulk-create

![](https://github.com/ntsk/gh-issue-bulk-create/actions/workflows/ci.yml/badge.svg)

A GitHub CLI extension to create multiple GitHub issues in bulk.

## Installation

```
gh extension install ntsk/gh-issue-bulk-create
```

## Usage

This extension uses a template markdown file and a CSV file containing data to create multiple GitHub issues in bulk.

```
gh issue-bulk-create --template <template_file> --csv <csv_file> [--repo <owner/repo>] [--dry-run]
```

### Options

- `--template`: Path to the template markdown file (required)
- `--csv`: Path to the CSV file containing data (required)
- `--repo`: Target repository in the format of owner/repo (default: current repository)
- `--dry-run`: Only show the content of issues without creating them

### Template File

The template file follows the GitHub Issue template format with front matter metadata at the beginning of the file:

```markdown
---
title: "{{title}}"
labels: {{label1}}, {{label2}}
assignees: {{assignee}}
---

## Description
{{description}}

## Steps to Reproduce
{{steps}}
```

You can use mustache syntax (`{{variable_name}}`) in the template to embed data from the CSV file.

### CSV File

The CSV file should contain headers that match the variable names used in the template.
The first line is the header row, and each subsequent line will be used to create a separate issue.

```csv
title,label1,label2,assignee,description,steps
Login page error,bug,frontend,username,Error appears when clicking login button,Click login button
Search not working,bug,backend,username,"Results don't appear when searching,Enter ""test"" in search box and click search"
```

## Example

You can try it with the sample files included in the repository:

```
# From inside a git repository directory
gh issue-bulk-create --template sample-template.md --csv sample-data.csv --dry-run

# Or specify a different repository
gh issue-bulk-create --template sample-template.md --csv sample-data.csv --repo owner/repo-name
```
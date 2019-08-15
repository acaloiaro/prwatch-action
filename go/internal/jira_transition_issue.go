package internal

import (
	"fmt"
	"log"
	"net/url"
	"os"

	jira "github.com/andygrunwald/go-jira"
)

// TransitionIssue transitions an issue's status to the one specified by the NEW_ISSUE_STATUS environment variable.
func TransitionIssue(issue Issue) {

	transitionName := os.Getenv("CONFLICT_ISSUE_STATUS")
	if transitionName == "" {
		log.Fatal("please set CONFLICT_ISSUE_STATUS with the status for in-conflict PRs, e.g. 'In Progress'")
	}

	cl := jiraClient()
	trs, _, err := cl.Issue.GetTransitions(issue.ID)
	if err != nil {
		log.Fatalf("unable to retrieve possible transition list for issue %v: %v", issue.ID, err)
	}

	// Find the desired transitionName in the possible transition list for this issue
	transitionId := ""
	for _, tr := range trs {
		if tr.Name == transitionName {
			transitionId = tr.ID
		}
	}

	if transitionId == "" {
		log.Fatalf("%s is not a valid transition for issue: %s", transitionName, issue.ID)
	}
	_, err = cl.Issue.DoTransition(issue.ID, transitionId)
	if err != nil {
		log.Fatalf("unable to transition issue: %v", err)
	}
}

func jiraClient() *jira.Client {
	jiraUser := os.Getenv("JIRA_USER")
	if jiraUser == "" {
		log.Fatal("Please set JIRA_USER environment variable with your Jira username")
	}

	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Please set JIRA_API_TOKEN environment variable with your Jira API token.")
	}

	jiraUrl := os.Getenv("JIRA_HOST")
	if jiraUrl == "" {
		log.Fatal("Please set JIRA_HOST environment variable with your Jira instance's hostname.")
	}

	url := fmt.Sprintf("https://%s:%s@%s", url.QueryEscape(jiraUser), apiToken, jiraUrl)
	jiraClient, err := jira.NewClient(nil, url)

	if err != nil {
		log.Fatal("Unable to connect to Jira:", err)
	}

	return jiraClient
}

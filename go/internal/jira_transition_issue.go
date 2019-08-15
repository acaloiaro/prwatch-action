package internal

import (
	"fmt"
	"log"
	"net/url"
	"os"

	jira "github.com/andygrunwald/go-jira"
)

var transitionName = os.Getenv("CONFLICT_ISSUE_STATUS")

// TransitionIssue transitions an issue's status to the one specified by the NEW_ISSUE_STATUS environment variable.
func TransitionIssue(issue *jira.Issue) {

	if transitionName == "" {
		log.Fatal("please set CONFLICT_ISSUE_STATUS with the status for in-conflict PRs, e.g. 'In Progress'")
	}

	cl := jiraClient()
	trs, _, err := cl.Issue.GetTransitions(issue.ID)
	if err != nil {
		log.Fatalf("unable to retrieve possible transition list for issue %v: %v", issue.ID, err)
	}

	// Find the desired transitionName in the possible transition list for this issue
	transitionID := ""
	for _, tr := range trs {
		if tr.Name == transitionName {
			transitionID = tr.ID
		}
	}

	if transitionID == "" {
		log.Fatalf("%s is not a valid transition for issue: %s", transitionName, issue.ID)
	}

	issue, _, err = cl.Issue.Get(issue.ID, nil)

	if err != nil || !shouldTransition(issue, transitionName) {
		log.Printf("not transitioning issue: %s", issue.ID)
		return
	}

	_, err = cl.Issue.DoTransition(issue.ID, transitionID)
	if err != nil {
		log.Printf("unable to transition issue: %v", err)
	}

	_, _, err = cl.Issue.AddComment(issue.ID, genComment(issue, transitionName))
	if err != nil {
		log.Printf("unable to leave comment on issue '%s': %v", issue.ID, err)
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

	jiraHost := os.Getenv("JIRA_HOST")
	if jiraHost == "" {
		log.Fatal("Please set JIRA_HOST environment variable with your Jira instance's hostname.")
	}

	url := fmt.Sprintf("https://%s:%s@%s", url.QueryEscape(jiraUser), apiToken, jiraHost)
	jiraClient, err := jira.NewClient(nil, url)

	if err != nil {
		log.Fatal("Unable to connect to Jira:", err)
	}

	return jiraClient
}

func shouldTransition(issue *jira.Issue, newStatus string) bool {

	currentStatus := issue.Fields.Status.Name

	// TODO: Make this list less brittle
	if currentStatus == newStatus ||
		currentStatus == "Archived" ||
		currentStatus == "Done" ||
		currentStatus == "Released" ||
		currentStatus == "Backlog" {
		return false
	}

	log.Printf("transitioning issue '%s' from '%s' to '%s'", issue.Key, currentStatus, newStatus)

	return true
}

func genComment(issue *jira.Issue, newStatus string) *jira.Comment {
	return &jira.Comment{
		Body: fmt.Sprintf("[~%s]: This card (%s) has been sent back to '%s' because its Pull Request has a merge conflict.", issue.Fields.Assignee.Key, issue.Key, newStatus),
	}
}

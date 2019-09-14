package internal

import (
	"fmt"
	"log"
	"net/url"

	"github.com/acaloiaro/prwatch/internal/config"
	jira "github.com/andygrunwald/go-jira"
)

type jiraIssueProvider struct {
	c *jira.Client
}

func newJiraIssueProvider(c *jira.Client) issueProvider {
	j := jiraIssueProvider{
		c: c,
	}

	return &j
}

// CommentIssue comments on jira issues with a pre-defined comment
func (j *jiraIssueProvider) CommentIssue(i issue) (ok bool) {

	if !config.SettingEnabled(i.Owner, "issue_comments") {
		return
	}

	jiraIssue, _, err := j.c.Issue.Get(i.ID, nil)
	if err != nil {
		log.Printf("unable to retrieve issue: '%s': %v", i.ID, err)
		return
	}

	transitionName := config.GetString("issues", "conflict_status")
	_, _, err = j.c.Issue.AddComment(i.ID, j.genComment(jiraIssue, transitionName))
	if err != nil {
		log.Printf("unable to leave comment on issue: '%s': %v", i.ID, err)
	}

	ok = err == nil

	return
}

func (j *jiraIssueProvider) jiraIssueFor(i issue) *jira.Issue {
	return &jira.Issue{
		ID:   i.ID,
		Key:  i.Key,
		Self: i.Value,
	}
}

// TransitionIssue transitions an issue's status to the one specified by the NEW_ISSUE_STATUS environment variable.
func (j *jiraIssueProvider) TransitionIssue(i issue) (ok bool) {

	// TODO: Use a constant for this setting name
	if !config.SettingEnabled(i.Owner, "transition_issues") {
		return
	}

	transitionName := config.GetString("issues", "conflict_status")
	if transitionName == "" {
		log.Fatal("Please set in config.yaml: config.issues.conflict_status. e.g. 'In Progress'")
	}

	trs, _, err := j.c.Issue.GetTransitions(i.ID)
	if err != nil {
		log.Printf("unable to retrieve possible transition list for issue %v: %v", i.ID, err)
		return
	}

	// Find the desired transitionName in the possible transition list for this issue
	transitionID := ""
	for _, tr := range trs {
		if tr.Name == transitionName {
			transitionID = tr.ID
		}
	}

	if transitionID == "" {
		log.Printf("%s is not a valid transition for issue: %s", transitionName, i.ID)
		return
	}

	jiraIssue, _, err := j.c.Issue.Get(i.ID, nil)
	if err != nil || !j.shouldTransition(jiraIssue, transitionName) {
		log.Printf("Not transitioning issue: %s.", i.ID)
		return
	}

	_, err = j.c.Issue.DoTransition(i.ID, transitionID)
	if err != nil {
		log.Printf("unable to transition issue: %v", err)
	}

	ok = true
	return
}

func newJiraClient() *jira.Client {
	if !config.SettingEnabled("jira") {
		log.Fatal("Jira must be enabled in config.yaml to create a jira client")
	}

	jiraUser := config.GetString("jira", "user")
	if jiraUser == "" {
		log.Fatal("Please set in config.yaml: settings.jira.user")
	}

	apiToken := config.GetEnv("JIRA_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Please set JIRA_API_TOKEN environment variable with your Jira API token.")
	}

	jiraHost := config.GetString("jira", "host")
	if jiraHost == "" {
		log.Fatal("Please set in config.yaml: settings.jira.host")
	}

	url := fmt.Sprintf("https://%s:%s@%s", url.QueryEscape(jiraUser), apiToken, jiraHost)
	jiraClient, err := jira.NewClient(nil, url)

	if err != nil {
		log.Fatal("Unable to connect to Jira:", err)
	}

	return jiraClient
}

func (j *jiraIssueProvider) shouldTransition(issue *jira.Issue, newStatus string) bool {

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

func (p *jiraIssueProvider) genComment(issue *jira.Issue, newStatus string) *jira.Comment {
	return &jira.Comment{
		Body: fmt.Sprintf("[~%s]: This card (%s) has been sent back to '%s' because its Pull Request has a merge conflict.", issue.Fields.Assignee.Key, issue.Key, newStatus),
	}
}

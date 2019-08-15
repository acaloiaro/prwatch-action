package main

import (
	"log"

	"github.com/acaloiaro/prwatch/internal"
	jira "github.com/andygrunwald/go-jira"
)

func main() {

	log.Println("Running...")
	githubClient := internal.NewClient()
	pulls, err := internal.ListPulls(githubClient)
	if err != nil {
		log.Println("Unable to fetch pull requests for repository: ", err)
	}

	var issueID string
	var ok bool

	for _, pull := range pulls {

		if issueID, ok = internal.IssueID(pull); !ok {
			continue
		}

		if internal.HasConflict(pull) {
			internal.TransitionIssue(&jira.Issue{ID: issueID})
		}
	}

	log.Println("Finished...")
}

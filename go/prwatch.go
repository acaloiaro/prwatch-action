package main

import (
	"log"

	"github.com/acaloiaro/prwatch/internal"
	"github.com/andygrunwald/go-jira"
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

		log.Println("Checking pull request:", pull.Number)

		if issueID, ok = internal.IssueID(pull); !ok {
			continue
		}

		if internal.HasConflict(pull) {

			log.Println("Pull request has conflict:", pull.Number)

			internal.TransitionIssue(&jira.Issue{ID: issueID})
		}
	}

	log.Println("Finished...")
}

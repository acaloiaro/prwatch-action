package main

import (
	"log"

	"github.com/acaloiaro/prwatch/internal"
)

func main() {

	githubClient := internal.NewClient()
	pulls, err := internal.ListPulls(githubClient)
	if err != nil {
		log.Println("Unable to fetch pull requests for repository: ", err)
	}

	for _, pull := range pulls {
		issueId := internal.IssueId(pull)
		if internal.HasConflict(pull) && issueId != "" {
			internal.TransitionIssue(internal.Issue{ID: issueId})
		}
	}
}

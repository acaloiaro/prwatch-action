package main

import (
	"github.com/acaloiaro/prwatch/internal"
	jira "github.com/andygrunwald/go-jira"
)

func main() {

	internal.TransitionIssue(&jira.Issue{ID: "GREEN-18854"})
}

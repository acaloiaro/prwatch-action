package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const defaultPageSize = 100

type actor struct {
	Login githubv4.String
}

type GithubPullRequest struct {
	Number    githubv4.Int
	Title     githubv4.String
	BodyText  githubv4.String
	Author    actor
	URL       githubv4.String
	UpdatedAt githubv4.DateTime
	Mergeable githubv4.MergeableState
}

type gqlRepository struct {
	Name  githubv4.String
	Owner owner
}

type owner struct {
	Login string
}

type pageInfo struct {
	EndCursor   githubv4.String
	HasNextPage githubv4.Boolean
}

type pullRequests struct {
	Nodes    []GithubPullRequest
	PageInfo pageInfo
}

type pullRequestQuery struct {
	Repository struct {
		PullRequests pullRequests `graphql:"pullRequests(states: [OPEN], first: $pageSize, orderBy: {field: UPDATED_AT, direction: ASC}, after: $pullsCursor)"`
	} `graphql:"repository(owner: $owner, name: $repository)"`
}

// ListPulls lists all open pulls requests for the current repository
func ListPulls(client GithubQueryer) (pulls []GithubPullRequest, err error) {
	o, repository, err := repositoryDetails()
	if err != nil {
		return
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(o),
		"repository":  githubv4.String(repository),
		"pullsCursor": (*githubv4.String)(nil),
		"pageSize":    githubv4.Int(defaultPageSize),
	}

	var query pullRequestQuery
	for {

		err = client.Query(&query, variables)
		if err != nil {
			return
		}

		pulls = append(pulls, query.Repository.PullRequests.Nodes...)

		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["pullsCursor"] = githubv4.NewString(query.Repository.PullRequests.PageInfo.EndCursor)
	}

	return
}

// HasConclift determines whether a pull request has a merge conflict
func HasConflict(pr GithubPullRequest) bool {
	return pr.Mergeable == "CONFLICTING"
}

func IssueID(pr GithubPullRequest) (issueID string, ok bool) {
	if len(string(pr.BodyText)) == 0 {
		ok = false
		return
	}

	// TODO: Make project-issue pattern more configurable
	re := regexp.MustCompile(fmt.Sprintf("%s-\\d*", os.Getenv("JIRA_PROJECT_NAME")))
	issueID = re.FindString(string(pr.BodyText))
	ok = issueID != ""
	return
}

func repositoryDetails() (owner, repository string, err error) {
	repoDetails := os.Getenv("GITHUB_REPOSITORY")
	details := strings.Split(repoDetails, "/")
	if len(details) != 2 {
		err = errors.New("Unable to determine the owner and repository where this Action is running. Check GITHUB_REPOSITORY")
		return
	}

	owner = details[0]
	repository = details[1]

	log.Printf("repo owner: '%s' repo name: '%s'", owner, repository)
	return
}

// GithubQueryer is an interface for performing github v4 graphql queries
type GithubQueryer interface {
	Query(query interface{}, variables map[string]interface{}) error
}

type githubClient struct {
	v4Client *githubv4.Client
	ctx      context.Context
}

// NewClient creates a new Github client
func NewClient() (client GithubQueryer) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	ctx := context.Background()
	return &githubClient{
		v4Client: githubv4.NewClient(oauth2.NewClient(ctx, src)),
		ctx:      ctx,
	}
}

// Query queries the github v4 graphql API
func (c *githubClient) Query(query interface{}, variables map[string]interface{}) error {
	return c.v4Client.Query(c.ctx, query, variables)
}

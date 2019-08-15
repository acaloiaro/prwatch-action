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

// GithubClient wraps the v4 graphql client to provide a higher level API
type GithubClient struct {
	v4Client *githubv4.Client
	ctx      context.Context
}

const defaultPageSize = 100

type actor struct {
	Login githubv4.String
}

type githubPullRequest struct {
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
	Nodes    []githubPullRequest
	PageInfo pageInfo
}

func NewClient() (client *GithubClient) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	ctx := context.Background()
	return &GithubClient{
		v4Client: githubv4.NewClient(oauth2.NewClient(ctx, src)),
		ctx:      ctx,
	}
}

func ListPulls(client *GithubClient) (pulls []githubPullRequest, err error) {
	if err != nil {
		log.Panic(err)
	}

	pulls, err = fetchPulls(client)

	return
}

// HasConclift determines whether a pull request has a merge conflict
func HasConflict(pr githubPullRequest) bool {
	return pr.Mergeable == "CONFLICTING"
}

func IssueId(pr githubPullRequest) (id string) {
	if len(string(pr.BodyText)) == 0 {
		return
	}

	re := regexp.MustCompile(fmt.Sprintf("%s-\\d*", os.Getenv("JIRA_PROJECT_NAME")))
	return re.FindString(string(pr.BodyText))
}

func fetchPulls(client *GithubClient) (pulls []githubPullRequest, err error) {

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

	var query struct {
		Repository struct {
			PullRequests pullRequests `graphql:"pullRequests(states: [OPEN], first: $pageSize, orderBy: {field: UPDATED_AT, direction: ASC}, after: $pullsCursor)"`
		} `graphql:"repository(owner: $owner, name: $repository)"`
	}

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

func repositoryDetails() (owner, repository string, err error) {
	repoDetails := os.Getenv("GITHUB_REPOSITORY")
	details := strings.Split(repoDetails, "/")
	if len(details) != 2 {
		err = errors.New("Unable to determine the owner and repository where this Action is running. Check GITHUB_REPOSITORY")
		return
	}

	owner = details[0]
	repository = details[1]
	return
}

func (c *GithubClient) Query(query interface{}, variables map[string]interface{}) error {
	return c.v4Client.Query(c.ctx, query, variables)
}

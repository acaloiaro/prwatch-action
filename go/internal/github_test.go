package internal

import (
	"errors"
	"os"
	"testing"

	"github.com/shurcooL/githubv4"
)

func TestRepoDetails(t *testing.T) {

	os.Setenv("GITHUB_REPOSITORY", "foo/bar")

	owner, repo, err := repositoryDetails()

	if err != nil {
		t.Error(err)
	}

	if owner != "foo" {
		t.Error("Owner should be 'foo'")
	}

	if repo != "bar" {
		t.Error("Repository should be 'bar'")
	}

	os.Setenv("GITHUB_REPOSITORY", "invalid")
	_, _, err = repositoryDetails()
	if err == nil {
		t.Error("Repository details should have returned an error")
	}
}

type MockClient struct {
	f func(query interface{}, variables map[string]interface{}) error
}

func (c *MockClient) Query(query interface{}, variables map[string]interface{}) error {
	return c.f(query, variables)
}

func TestListPulls(t *testing.T) {

	os.Setenv("GITHUB_REPOSITORY", "acaloiaro/isok")
	client := &MockClient{}

	goodQuery := func(query interface{}, v map[string]interface{}) error {
		q := query.(*pullRequestQuery)

		q.Repository.PullRequests = pullRequests{Nodes: []GithubPullRequest{GithubPullRequest{Number: 1}}}
		query = q

		expectedOwner := "acaloiaro"
		if string(v["owner"].(githubv4.String)) != expectedOwner {
			t.Errorf("expected repository to be: '%s', got: '%s'", expectedOwner, v["owner"])
		}

		expectedRepo := "isok"
		if string(v["repository"].(githubv4.String)) != expectedRepo {
			t.Errorf("expected repository to be: '%s', got: '%s'", expectedRepo, v["repository"])
		}

		return nil
	}
	client.f = goodQuery

	pulls, err := ListPulls(client)

	if err != nil {
		t.Error(err)
	}

	firstPull := pulls[0]
	if firstPull.Number != 1 {
		t.Error("first PR number should have been 1")
	}

	badQuery := func(query interface{}, v map[string]interface{}) error {
		return errors.New("bad things happened")
	}
	client.f = badQuery

	_, err = ListPulls(client)
	if err == nil {
		t.Error("should get an error when the client fails")
	}

	// TODO: Add tests over pagination of results 
}

func TestIssueId(t *testing.T) {

	os.Setenv("JIRA_PROJECT_NAME", "FOO")

	pr := GithubPullRequest{
		BodyText: "Issue url is https://foobar.atlassian.net/browse/FOO-1234",
	}

	const expectedID = "FOO-1234"
	if ID, ok := IssueID(pr); ID != expectedID || !ok {
		t.Errorf("expected issue id: %s: got: %s", expectedID, ID)
	}
}

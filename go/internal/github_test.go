package internal

import (
	"errors"
	"testing"

	"github.com/acaloiaro/prwatch/internal/config"
	"github.com/shurcooL/githubv4"
)

type mockFilesProvider struct {
	files map[string]bool
}

func (p mockFilesProvider) Exists(path string) bool {
	return p.files[path]
}

func TestRepoDetails(t *testing.T) {

	config.SetEnv("GITHUB_REPOSITORY", "foo/bar")

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

	config.SetEnv("GITHUB_REPOSITORY", "invalid")
	_, _, err = repositoryDetails()
	if err == nil {
		t.Error("Repository details should have returned an error")
	}
}

type MockGithubClient struct {
	f         func(query interface{}, variables map[string]interface{}) error
	pageCount int
}

func (c *MockGithubClient) Query(query interface{}, variables map[string]interface{}) error {
	return c.f(query, variables)
}

func TestListPulls(t *testing.T) {

	config.SetEnv("GITHUB_REPOSITORY", "acaloiaro/isok")
	client := &MockGithubClient{}

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

	// pagination test
	expectedPages := 5
	client.pageCount = 1
	client.f = func(query interface{}, v map[string]interface{}) error {
		morePages := false
		if client.pageCount < 5 {
			morePages = true
		}

		q := query.(*pullRequestQuery)
		q.Repository.PullRequests = pullRequests{
			PageInfo: pageInfo{HasNextPage: githubv4.Boolean(morePages)},
			Nodes:    []GithubPullRequest{GithubPullRequest{Number: githubv4.Int(client.pageCount)}},
		}

		client.pageCount = client.pageCount + 1
		return nil
	}

	pulls, err = ListPulls(client)
	numPulls := len(pulls)
	if numPulls != expectedPages {
		t.Errorf("expected to have paged results from client. expected: '%d' results, got: '%d'", expectedPages, numPulls)
	}

}

func TestIssueId(t *testing.T) {

	config.GlobalEnable(config.Jira)
	config.GlobalSet(config.JiraProjectName, "FOO")

	pr := GithubPullRequest{
		BodyText: "Issue url is https://foobar.atlassian.net/browse/FOO-1234",
	}

	const expectedID = "FOO-1234"
	if ID, ok := IssueID(pr); ID != expectedID || !ok {
		t.Errorf("expected issue id: %s: got: %s", expectedID, ID)
	}
}

func TestHasConflict(t *testing.T) {

	defer services.reset()

	pr := GithubPullRequest{
		Mergeable: githubv4.MergeableStateConflicting,
	}

	// when .gitattributes doesn't exist and the PR is in conflict, then there is a conflict
	services.f = mockFilesProvider{files: map[string]bool{".gitattributes": false}}
	services.g = &mockGitProvider{}
	conflict := hasConflict(pr)
	if !conflict {
		t.Error("this pull request should be considered in conflict")
	}

	// when .gitattributes exists and the PR is in conflict, then there is a conflict only when merging fails
	services.f = mockFilesProvider{files: map[string]bool{".gitattributes": true}}
	services.g = &mockGitProvider{mergeFunc: func(ref string, a ...string) error { return errors.New("no good") }}
	conflict = hasConflict(pr)
	if !conflict {
		t.Error("this pull request should be considered in conflict")
	}

	// when .gitattributes exists and the PR is in conflict, then there is a conflict only when merging fails
	services.f = mockFilesProvider{files: map[string]bool{".gitattributes": true}}
	services.g = &mockGitProvider{mergeFunc: func(ref string, a ...string) error { return nil }}
	conflict = hasConflict(pr)
	if conflict {
		t.Error("this pull request should not be considered in conflict")
	}
}

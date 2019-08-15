package internal

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Error("New client creation failed")
	}

	if client.ctx == nil {
		t.Error("Autentication should have initialized a context")
	}

}

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

func TestListPulls(t *testing.T) {

	os.Setenv("GITHUB_REPOSITORY", "grnhse/jben")

	client := NewClient()
	pulls, err := ListPulls(client)

	if err != nil {
		t.Error(err)
	}

	/*for _, pull := range pulls {*/
	//t.Error(pull)
	/*}*/
	if len(pulls) == 0 {
		t.Error("Should have listed some PRs")
	}
}

func TestIssueId(t *testing.T) {
	pr := githubPullRequest{
		BodyText: "Issue url is https://foobar.atlassian.net/browse/GREEN-1234",
	}

	if IssueId(pr) != "GREEN-1234" {
		t.Errorf("Did not match: %s", IssueId(pr))
	}
}

package internal

import (
	"errors"
	"testing"
)

type mockGitProvider struct {
	checkoutFunc   func(ref string) error
	currentRefFunc func() string
	mergeFunc      func(ref string, args ...string) error
	resetFunc      func(ref string, args ...string) error
}

func (e *mockGitProvider) CurrentRefName() string {
	if e.currentRefFunc != nil {
		return e.currentRefFunc()
	}

	return "undefined"
}

func (e *mockGitProvider) Checkout(ref string) error {
	if e.checkoutFunc != nil {
		return e.checkoutFunc(ref)
	}

	return nil
}

func (e *mockGitProvider) Merge(ref string, args ...string) error {
	if e.mergeFunc != nil {
		return e.mergeFunc(ref, args...)
	}

	return nil
}

func (e *mockGitProvider) Reset(ref string, args ...string) error {
	if e.resetFunc != nil {
		return e.resetFunc(ref, args...)
	}

	return nil
}

func TestTryMerge(t *testing.T) {

	// leave services in a good state for other tests
	defer services.reset()

	pr := GithubPullRequest{
		BaseRefName: "foo",
		HeadRefName: "bar",
	}

	services.git = &mockGitProvider{checkoutFunc: func(ref string) error { return errors.New("fail") }}
	status := tryMerge(pr)
	if status {
		t.Error("Should not have been able to merge")
	}

	services.git = &mockGitProvider{mergeFunc: func(ref string, a ...string) error { return errors.New("fail") }}
	status = tryMerge(pr)
	if status {
		t.Error("Should not have been able to merge")
	}

	services.git = &mockGitProvider{resetFunc: func(ref string, a ...string) error { return errors.New("fail") }}
	status = tryMerge(pr)
	if status {
		t.Error("Should not have been able to merge")
	}

	services.git = &mockGitProvider{}
	status = tryMerge(pr)
	if !status {
		t.Error("Should have been able to merge")
	}
}

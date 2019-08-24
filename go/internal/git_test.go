package internal

import (
	"errors"
	"testing"
)

type MockGitExecutor struct {
	checkoutFunc   func(ref string) error
	currentRefFunc func() string
	mergeFunc      func(ref string, args ...string) error
	resetFunc      func(ref string, args ...string) error
}

func (e *MockGitExecutor) CurrentRefName() string {
	if e.currentRefFunc != nil {
		return e.currentRefFunc()
	}

	return "undefined"
}

func (e *MockGitExecutor) Checkout(ref string) error {
	if e.checkoutFunc != nil {
		return e.checkoutFunc(ref)
	}

	return nil
}

func (e *MockGitExecutor) Merge(ref string, args ...string) error {
	if e.mergeFunc != nil {
		return e.mergeFunc(ref, args...)
	}

	return nil
}

func (e *MockGitExecutor) Reset(ref string, args ...string) error {
	if e.resetFunc != nil {
		return e.resetFunc(ref, args...)
	}

	return nil
}

func TestTryMerge(t *testing.T) {

	pr := GithubPullRequest{
		BaseRefName: "foo",
		HeadRefName: "bar",
	}

	status := TryMerge(&MockGitExecutor{checkoutFunc: func(ref string) error { return errors.New("fail") }}, pr)
	if status {
		t.Error("Should not have been able to merge")
	}

	status = TryMerge(&MockGitExecutor{mergeFunc: func(ref string, a ...string) error { return errors.New("fail") }}, pr)
	if status {
		t.Error("Should not have been able to merge")
	}

	status = TryMerge(&MockGitExecutor{resetFunc: func(ref string, a ...string) error { return errors.New("fail") }}, pr)
	if status {
		t.Error("Should not have been able to merge")
	}

	status = TryMerge(&MockGitExecutor{}, pr)
	if !status {
		t.Error("Should have been able to merge")
	}
}

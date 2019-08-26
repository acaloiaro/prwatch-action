package internal

import (
	"errors"
	"testing"
	"time"
)

type mockGitProvider struct {
	checkoutCalled   time.Time
	checkoutFunc     func(ref string) error
	currentRefCalled time.Time
	currentRefFunc   func() string
	mergeCalled      time.Time
	mergeFunc        func(ref string, args ...string) error
	resetCalled      time.Time
	resetFunc        func(ref string, args ...string) error
}

func (e *mockGitProvider) CurrentRefName() string {

	e.currentRefCalled = time.Now()

	if e.currentRefFunc != nil {
		return e.currentRefFunc()
	}

	return "undefined"
}

func (e *mockGitProvider) Checkout(ref string) error {

	e.checkoutCalled = time.Now()

	if e.checkoutFunc != nil {
		return e.checkoutFunc(ref)
	}

	return nil
}

func (e *mockGitProvider) Merge(ref string, args ...string) error {

	e.mergeCalled = time.Now()

	if e.mergeFunc != nil {
		return e.mergeFunc(ref, args...)
	}

	return nil
}

func (e *mockGitProvider) Reset(ref string, args ...string) error {

	e.resetCalled = time.Now()

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

	p := &mockGitProvider{}
	services.git = p
	status = tryMerge(pr)
	if !status {
		t.Error("Should have been able to merge")
	}

	if p.currentRefCalled.Before(p.checkoutCalled) || p.mergeCalled.Before(p.checkoutCalled) || p.resetCalled.Before(p.mergeCalled) {
		t.Error("git operations called in the wrong order")
	}
}

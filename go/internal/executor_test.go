package internal

import (
	"testing"
	"time"

	"github.com/acaloiaro/prwatch/internal/config"
)

type testExecutionPlan struct {
	passCount          uint
	firstPassFinished  time.Time
	secondPassFinished time.Time
	githubClient       GithubQueryer
	t                  *time.Timer
	f                  func() error
}

func (t *testExecutionPlan) Execute() (err error) {

	err = t.f()

	return
}

func (t testExecutionPlan) DualPassTimer() *time.Timer {
	return t.t
}

func (t testExecutionPlan) client() GithubQueryer {
	return t.githubClient
}

func TestExecutorExecute(t *testing.T) {
	client := &MockGithubClient{}

	config.SetEnv("GITHUB_REPOSITORY", "acaloiaro/isok")

	config.GlobalEnable("dual_pass")
	config.GlobalSet("dual_pass", "wait_duration", "1ms")

	st := &testExecutionPlan{
		passCount:    0,
		githubClient: client,
		t:            time.NewTimer(dualPassInterval()), // DualPassTimer
	}

	// The second pass occurrs when the execution strategy's Execute() is called; f implements the strategy's Execute()
	secondPass := func() error {
		st.secondPassFinished = time.Now()
		st.passCount = st.passCount + 1
		return nil
	}
	st.f = secondPass

	// The first "pass" occurrs when the first ListPulls call is issued to the github client
	firstPass := func(query interface{}, v map[string]interface{}) error {
		st.passCount = 1
		st.firstPassFinished = time.Now()
		return nil
	}
	client.f = firstPass

	e := NewExecutor(st)
	err := e.Execute()

	if err != nil {
		t.Error(err)
	}

	if st.passCount != 2 {
		t.Errorf("should have been two passes: %v", st.passCount)
	}

	if !st.secondPassFinished.After(st.firstPassFinished) {
		t.Error("phase 1 should have finished before phase 2")
	}
}

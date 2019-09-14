// Package internal executors are the entrypoint to prwatch. They define the overall behavior of the action.
package internal

import (
	"log"
	"math"
	"time"

	"github.com/acaloiaro/prwatch/internal/config"
)

type executor struct {
	executionPlan executionPlan
}

// NewExecutor creates a new executor for the provided execution plan
func NewExecutor(s executionPlan) *executor {
	return &executor{
		executionPlan: s,
	}
}

// Execute executes the executionPlan in either one or two phases.
// Dual phase mode is designed to retrieve Github's mergability status in two phases when this action is triggered on
// merge to a base branch. Because immediately following a merge, Github cannot yet determine the mergability of pull
// requests, the first phase is a request to Github to update its mergability statuses.
//
// The second phase is to determine the actual mergability of all open pull requests.
func (e *executor) Execute() error {

	if config.SettingEnabled(config.DualPass) {
		timer := e.executionPlan.DualPassTimer()

		// List open pull requests to trigger a refresh of Github's mergability status
		ListPulls(e.executionPlan.client())

		done := false
		tick := time.Tick(1 * time.Second)
		countdown := time.Now().Add(dualPassInterval())

		for !done {
			select {
			case <-tick:
				s := math.Round(countdown.Sub(time.Now()).Seconds())
				log.Println("Waiting ...", s, "seconds")
			case <-timer.C:
				log.Println("Phase 1 complete.")
				done = true
			}
		}
	} else {
		log.Println("Single pass mode")
	}

	return e.executionPlan.Execute()
}

type executionPlan interface {
	Execute() error
	DualPassTimer() *time.Timer
	client() GithubQueryer
}

// DefaultExecutionPlan is the executionPlan used by the main executable
type DefaultExecutionPlan struct {
	GithubClient GithubQueryer
}

// Execute executes an executionPlan
func (e *DefaultExecutionPlan) Execute() error {
	pulls, err := ListPulls(e.GithubClient)
	if err != nil {
		log.Println("Unable to fetch pull requests for repository: ", err)
		return err
	}

	var issueID string
	var ok bool

	for _, pull := range pulls {

		log.Println("checking pull request:", pull.Number)

		if issueID, ok = IssueID(pull); !ok {
			log.Printf("no issue ID associated with this pull request '%d', skipping", pull.Number)
			continue
		}

		if hasConflict(pull) {

			log.Printf("pull request has conflict: %s", pull.URL)

			i := issue{ID: issueID, Owner: string(pull.Author.Login)}

			services.issues().TransitionIssue(i)
			services.issues().CommentIssue(i)
		} else {
			log.Printf("pull request is not conflicitng: %s", pull.URL)
		}
	}

	return nil
}

func (e DefaultExecutionPlan) client() GithubQueryer {
	return e.GithubClient
}

// DualPassTimer returns a timer if DUAL_PASS_WAIT_DURATION contains a valid duration string, nil otherwise
func (e DefaultExecutionPlan) DualPassTimer() (timer *time.Timer) {

	if d := dualPassInterval(); d > 0 {
		timer = time.NewTimer(d)
	}

	return
}

func dualPassInterval() time.Duration {

	d := config.GetDuration(config.DualPassWaitDuration)

	return d
}

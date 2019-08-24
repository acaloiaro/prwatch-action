package internal

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// TryMerge attempts to merge a pull request locally
// The purpose of merging locally is because the Mergable status from Github is insufficient when a .gitatributes file
// is present. Because Github does not support custom merge drivers, e.g. `merge=union` from .gitattributes, merging
// using a git client that does support custom merge drivers is the only way to tell whether a branch is truly mergable.
func TryMerge(g GitExecutor, pr GithubPullRequest) bool {

	origBranchRef := g.CurrentRefName()

	// set the the current branch
	err := g.Checkout(string(pr.BaseRefName))
	if err != nil {
		return false
	}

	err = g.Merge(fmt.Sprintf("origin/%s", string(pr.HeadRefName)))
	if err != nil {
		return false
	}

	//reset HEAD back to the HEAD prior to merging
	err = g.Reset(origBranchRef, "--hard")

	return err == nil
}

// GitExecutor is an interface for performing various Git operations
type GitExecutor interface {
	Checkout(ref string) error
	CurrentRefName() string
	Merge(ref string, args ...string) error
	Reset(ref string, args ...string) error
}

// GitCommandLine is a GitExecutor for the command-line executable of git, i.e. "git" proper
// Until there is a 100% Golang git implementation that supports .gitattributes files, this is the preferred method of
// executing git operations on the local git repository.
type GitCommandLine struct{}

// CurrentRefName returns the ref of the current HEAD
func (gcl *GitCommandLine) CurrentRefName() string {

	// this Github action uses `actions/checkout`, which places the repo in a "detached head" state,
	// so this command gives us the sha of the detatched head
	origBranchRef, _ := exec.Command("git", "rev-parse", "HEAD").Output()

	return strings.TrimSpace(string(origBranchRef))
}

// Checkout checks out git reference `ref`
func (gcl *GitCommandLine) Checkout(ref string) error {

	out, err := exec.Command("git", "checkout", ref).CombinedOutput()
	if err != nil {
		log.Println("Error checking out branch:", string(out))
	}

	return err
}

// Merge merges git reference `ref` with the current HEAD, passing `args` to the merge command
func (gcl *GitCommandLine) Merge(ref string, args ...string) error {

	out, err := exec.Command(
		"git",
		"-c",
		"user.name=prwatch",
		"-c",
		"user.email=prwatch@github.bot",
		"merge",
		ref,
		"-m",
		"Test merge").CombinedOutput()

	if err != nil {
		log.Println("Error merging branch:", string(out))
	}

	return err
}

// Reset resets HEAD to git reference `ref`, passing `args` to the reset command
func (gcl *GitCommandLine) Reset(ref string, args ...string) error {

	combinedArgs := append([]string{"reset", string(ref)}, args...)

	err := exec.Command("git", combinedArgs...).Run()
	if err != nil {
		log.Println("Error resetting branch:", err, combinedArgs)
	}

	return err
}

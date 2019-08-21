package internal

import (
	"log"
	"os/exec"
)

func AttemptMerge(pr GithubPullRequest) bool {
	err := exec.Command("git", "checkout", string(pr.BaseRefName)).Run()
	if err != nil {
		log.Println("Unable to checkout base ref:", err)
		return false
	}

	cmd := exec.Command("git", "merge", string(pr.HeadRefName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Unable to merge branch:", output)
		return false
	}

	// Return the working directory to a clean state
	err = exec.Command("git", "reset", "origin/master", "--hard").Run()
	if err != nil {
		log.Println("Unable to reset branch to master (this will cause problems)")
	}

	return true
}

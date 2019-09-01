package main

import (
	"log"

	"github.com/acaloiaro/prwatch/internal"
)

func main() {

	log.Println("Running...")

	executor := internal.NewExecutor(&internal.DefaultExecutionPlan{GithubClient: internal.NewGithubClient()})
	err := executor.Execute()

	if err != nil {
		log.Printf("Finished unsuccessfully: %s", err)
		return
	}

	log.Println("Finished...")
}

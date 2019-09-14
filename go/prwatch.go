package main

import (
	"log"

	"github.com/acaloiaro/prwatch/internal"
	"github.com/acaloiaro/prwatch/internal/config"
)

func main() {

	log.Println("Running...")

	config.Initialize()
	executor := internal.NewExecutor(&internal.DefaultExecutionPlan{GithubClient: internal.NewGithubClient()})
	err := executor.Execute()

	if err != nil {
		log.Printf("Finished unsuccessfully: %s", err)
		return
	}

	log.Println("Finished...")
}

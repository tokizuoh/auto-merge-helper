package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// ref: https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	token := os.Getenv("GITHUB_TOKEN")
	repository := os.Getenv("GITHUB_REPOSITORY")
	sha := os.Getenv("GITHUB_SHA")
	if token == "" || repository == "" || sha == "" {
		log.Fatalln("failed to retrieve expected environment variable")
	}

	ownerRepoName := strings.Split(repository, "/")
	if len(ownerRepoName) != 2 {
		log.Fatalf("failed to split the \"%s\" as expected.", repository)
	}
	owner := ownerRepoName[0]
	repoName := ownerRepoName[1]

	query, err := fetch(token, owner, repoName, sha)
	if err != nil {
		log.Fatalln(err)
	}

	ok, err := checkAllSuccess(query.Repository.Object.Commit.StatusCheckRollup.Contexts.Nodes)
	if err != nil {
		log.Fatalln(err)
	}

	if ok {
		fmt.Println("All CI's are SUCCESS!")
	} else {
		log.Fatalln("failed to run CI successfully, at least one CI has failed")
	}
}
